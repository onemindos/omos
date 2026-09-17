"""
Printer Agent - 3D Printer Discovery, Slicing, and Print Job Management
Extracted from ADA V2 (https://github.com/nazirlouis/ada_v2)

This agent handles:
- Network printer discovery via mDNS (Zeroconf)
- OrcaSlicer/PrusaSlicer CLI integration for slicing
- OctoPrint and Moonraker API for print job submission
- Auto-detection of suitable slicing profiles

Dependencies:
    pip install zeroconf aiohttp

External Dependencies:
    - OrcaSlicer or PrusaSlicer installed on system

For OneMind Integration:
    - Add to hardware/ directory or backend/agno/tools/
    - Use NATS for status callbacks
    - Integrate with CAD agent for full design-to-print workflow
"""

import asyncio
import os
import subprocess
import json
import platform
import glob
import shutil
from typing import Dict, List, Optional, Any
from dataclasses import dataclass, asdict
from enum import Enum

import aiohttp
from zeroconf import Zeroconf, ServiceBrowser, ServiceListener


class PrinterType(Enum):
    """Supported 3D printer types."""

    OCTOPRINT = "octoprint"
    MOONRAKER = "moonraker"
    PRUSALINK = "prusalink"
    UNKNOWN = "unknown"


@dataclass
class Printer:
    """Represents a discovered 3D printer."""

    name: str
    host: str
    port: int
    printer_type: PrinterType
    api_key: Optional[str] = None
    camera_url: Optional[str] = None

    def to_dict(self) -> dict:
        d = asdict(self)
        d["printer_type"] = self.printer_type.value
        return d


@dataclass
class PrintStatus:
    """Current status of a print job."""

    printer: str
    state: str  # "printing", "idle", "paused", "error"
    progress_percent: float
    time_remaining: Optional[str]
    time_elapsed: Optional[str]
    filename: Optional[str]
    temperatures: Optional[Dict[str, Dict[str, float]]] = None

    def to_dict(self) -> dict:
        return asdict(self)


class PrinterDiscoveryListener(ServiceListener):
    """mDNS listener for printer discovery."""

    def __init__(self):
        self.printers: List[Printer] = []

    def add_service(self, zc: Zeroconf, type_: str, name: str) -> None:
        info = zc.get_service_info(type_, name)
        if info:
            host = info.parsed_addresses()[0] if info.parsed_addresses() else None
            # Fallback to server name if address parsing fails
            if not host and info.server:
                host = info.server.rstrip(".")

            if host:
                # Determine printer type from service type
                if "_octoprint._tcp" in type_:
                    printer_type = PrinterType.OCTOPRINT
                elif "_moonraker._tcp" in type_ or "_klipper._tcp" in type_:
                    printer_type = PrinterType.MOONRAKER
                else:
                    printer_type = PrinterType.UNKNOWN

                printer = Printer(
                    name=name.replace(f".{type_}", ""),
                    host=host,
                    port=info.port or 80,
                    printer_type=printer_type,
                )
                self.printers.append(printer)
                print(
                    f"[PRINTER] Discovered: {printer.name} at {printer.host}:{printer.port} ({printer.printer_type.value})"
                )

    def remove_service(self, zc: Zeroconf, type_: str, name: str) -> None:
        pass

    def update_service(self, zc: Zeroconf, type_: str, name: str) -> None:
        pass


class PrinterAgent:
    """
    Handles 3D printer discovery, profile management, slicing, and print job submission.
    """

    def __init__(self, profiles_dir: str = "printer_profiles"):
        self.printers: Dict[str, Printer] = {}  # host -> Printer
        self.profiles_dir = profiles_dir
        self._zeroconf: Optional[Zeroconf] = None
        self._error_tracker = set()  # Track hosts with errors to prevent log spam

        # Detect slicer path and profiles directory
        self.slicer_path = self._detect_slicer_path()
        self._orca_profiles_dir = self._detect_orca_profiles_dir()

        # Ensure profiles directory exists
        os.makedirs(profiles_dir, exist_ok=True)

    def _detect_orca_profiles_dir(self) -> Optional[str]:
        """Detect OrcaSlicer profiles directory."""
        system = platform.system()

        if system == "Darwin":  # macOS
            base = os.path.expanduser("~/Library/Application Support/OrcaSlicer")
        elif system == "Windows":
            base = os.path.join(os.environ.get("APPDATA", ""), "OrcaSlicer")
        else:  # Linux
            base = os.path.expanduser("~/.config/OrcaSlicer")

        if os.path.isdir(base):
            print(f"[PRINTER] Found OrcaSlicer profiles at: {base}")
            return base

        return None

    def get_available_profiles(self) -> Dict[str, List[str]]:
        """
        Get all available OrcaSlicer profiles from the system folder.
        Returns dict with 'machines', 'processes', 'filaments' lists.
        """
        profiles = {"machines": [], "processes": [], "filaments": []}

        if not self._orca_profiles_dir:
            return profiles

        system_dir = os.path.join(self._orca_profiles_dir, "system")
        if not os.path.isdir(system_dir):
            return profiles

        # Scan all vendor folders (Creality, Custom, etc.)
        for vendor in os.listdir(system_dir):
            vendor_path = os.path.join(system_dir, vendor)
            if not os.path.isdir(vendor_path):
                continue

            # Machine profiles
            machine_dir = os.path.join(vendor_path, "machine")
            if os.path.isdir(machine_dir):
                for f in os.listdir(machine_dir):
                    if f.endswith(".json"):
                        profiles["machines"].append(f"system/{vendor}/machine/{f}")

            # Process profiles
            process_dir = os.path.join(vendor_path, "process")
            if os.path.isdir(process_dir):
                for f in os.listdir(process_dir):
                    if f.endswith(".json"):
                        profiles["processes"].append(f"system/{vendor}/process/{f}")

            # Filament profiles
            filament_dir = os.path.join(vendor_path, "filament")
            if os.path.isdir(filament_dir):
                for f in os.listdir(filament_dir):
                    if f.endswith(".json"):
                        profiles["filaments"].append(f"system/{vendor}/filament/{f}")

        return profiles

    def _find_matching_profile(
        self, printer_name: str, profile_type: str
    ) -> Optional[str]:
        """
        Find a matching profile for a printer by name.
        profile_type: 'machine', 'process', or 'filament'
        """
        if not self._orca_profiles_dir:
            return None

        # Normalize printer name for matching
        search_terms = printer_name.lower().split()

        # Common brand keywords to identify vendor folder
        vendor_map = {
            "creality": "Creality",
            "ender": "Creality",
            "cr-": "Creality",
            "k1": "Creality",
            "prusa": "Prusa",
            "bambu": "Bambu Lab",
            "anycubic": "Anycubic",
        }

        # Try to identify vendor
        vendor = None
        for term in search_terms:
            for key, val in vendor_map.items():
                if key in term:
                    vendor = val
                    break
            if vendor:
                break

        if not vendor:
            vendor = "Creality"  # Default fallback

        system_dir = os.path.join(self._orca_profiles_dir, "system", vendor)
        if not os.path.isdir(system_dir):
            print(f"[PRINTER] Vendor folder not found: {vendor}")
            return None

        target_dir = os.path.join(system_dir, profile_type)
        if not os.path.isdir(target_dir):
            return None

        # Score-based matching
        best_match = None
        best_score = 0

        for filename in os.listdir(target_dir):
            if not filename.endswith(".json"):
                continue

            name_lower = filename.lower()
            score = 0

            # Score based on matching search terms
            for term in search_terms:
                if term in name_lower:
                    score += 10
                    if profile_type == "machine":
                        idx = name_lower.find(term)
                        if idx >= 0:
                            after_idx = idx + len(term)
                            if after_idx < len(name_lower):
                                next_char = name_lower[after_idx]
                                if next_char.isalpha():
                                    score -= 8  # Variant like K1C
                                elif next_char in " .(-":
                                    score += 5  # Direct match

            # Bonus for "0.4 nozzle" (most common)
            if "0.4" in name_lower:
                score += 2

            # Bonus for "standard" or "optimal" process profiles
            if profile_type == "process":
                if "standard" in name_lower:
                    score += 5
                elif "optimal" in name_lower:
                    score += 3

            # Bonus for generic PLA filament
            if profile_type == "filament":
                if "pla" in name_lower and "generic" in name_lower:
                    score += 5
                    if "-cf" in name_lower or "-gf" in name_lower:
                        score -= 5  # Carbon fiber / glass fiber
                    if "silk" in name_lower or "matte" in name_lower:
                        score -= 2  # Specialty finishes

            if score > best_score:
                best_score = score
                best_match = os.path.join(target_dir, filename)

        if best_match:
            print(
                f"[PRINTER] Matched {profile_type} profile: {os.path.basename(best_match)} (score: {best_score})"
            )

        return best_match

    def get_profiles_for_printer(self, printer_name: str) -> Dict[str, Optional[str]]:
        """
        Auto-detect suitable profiles for a given printer name.
        Returns dict with 'machine', 'process', 'filament' paths.
        """
        return {
            "machine": self._find_matching_profile(printer_name, "machine"),
            "process": self._find_matching_profile(printer_name, "process"),
            "filament": self._find_matching_profile(printer_name, "filament"),
        }

    def _detect_slicer_path(self) -> Optional[str]:
        """Detect OrcaSlicer or PrusaSlicer installation path."""
        system = platform.system()

        paths = []
        if system == "Darwin":  # macOS
            paths = [
                "/Applications/OrcaSlicer.app/Contents/MacOS/OrcaSlicer",
                "/Applications/PrusaSlicer.app/Contents/MacOS/PrusaSlicer",
                "/Applications/Original Prusa Drivers/PrusaSlicer.app/Contents/MacOS/PrusaSlicer",
            ]
        elif system == "Windows":
            paths = [
                r"C:\Program Files\OrcaSlicer\orca-slicer-console.exe",
                r"C:\Program Files\OrcaSlicer\orca-slicer.exe",
                r"C:\Program Files\Prusa3D\PrusaSlicer\prusa-slicer-console.exe",
                r"C:\Program Files (x86)\Prusa3D\PrusaSlicer\prusa-slicer-console.exe",
                os.path.join(
                    os.environ.get("LOCALAPPDATA", ""),
                    "Programs",
                    "OrcaSlicer",
                    "orca-slicer-console.exe",
                ),
            ]
        else:  # Linux
            paths = [
                "/usr/bin/orca-slicer",
                os.path.expanduser("~/.local/bin/orca-slicer"),
                "/usr/bin/prusa-slicer",
                "/usr/local/bin/prusa-slicer",
            ]

        for path in paths:
            if os.path.exists(path):
                print(f"[PRINTER] Found Slicer at: {path}")
                return path

        # Try to find via which/where
        for binary in ["orca-slicer", "prusa-slicer", "prusa-slicer-console"]:
            try:
                result = subprocess.run(
                    ["which", binary] if system != "Windows" else ["where", binary],
                    capture_output=True,
                    text=True,
                )
                if result.returncode == 0 and result.stdout.strip():
                    path = result.stdout.strip().split("\n")[0]
                    print(f"[PRINTER] Found Slicer via PATH: {path}")
                    return path
            except Exception:
                pass

        print("[PRINTER] Warning: No Slicer (Orca/Prusa) found. Slicing will fail.")
        return None

    async def discover_printers(self, timeout: float = 5.0) -> List[Dict]:
        """
        Discovers 3D printers on the local network via mDNS.
        Returns list of discovered printers.
        """
        print(f"[PRINTER] Starting printer discovery (timeout: {timeout}s)...")

        self._zeroconf = Zeroconf()
        listener = PrinterDiscoveryListener()

        # Browse for common 3D printer services
        services = [
            "_octoprint._tcp.local.",
            "_moonraker._tcp.local.",
            "_klipper._tcp.local.",
            "_http._tcp.local.",  # Generic HTTP - fallback
        ]

        browsers = []
        for service in services:
            browser = ServiceBrowser(self._zeroconf, service, listener)
            browsers.append(browser)

        # Wait for discovery
        await asyncio.sleep(timeout)

        # Cleanup
        self._zeroconf.close()

        # Probe unknown printers to identify type
        for printer in listener.printers:
            if printer.printer_type == PrinterType.UNKNOWN:
                print(f"[PRINTER] Probing unknown printer: {printer.host}...")
                ptype = await self._probe_printer_type(printer.host, printer.port)
                if ptype != PrinterType.UNKNOWN:
                    printer.printer_type = ptype
                    print(f"[PRINTER] Identified {printer.name} as {ptype.value}")

        # Probe for cameras
        for printer in listener.printers:
            if not printer.camera_url:
                cam_url = await self._probe_camera(printer.host, printer.port)
                if cam_url:
                    printer.camera_url = cam_url

        # Store discovered printers
        for printer in listener.printers:
            self.printers[printer.host] = printer

        print(f"[PRINTER] Discovery complete. Found {len(self.printers)} printers.")
        return [p.to_dict() for p in self.printers.values()]

    async def _probe_printer_type(self, host: str, port: int) -> PrinterType:
        """Probe a host to check if it's running Moonraker or OctoPrint."""
        try:
            timeout = aiohttp.ClientTimeout(total=2.0, connect=1.0)
            async with aiohttp.ClientSession(timeout=timeout) as session:
                # Check Moonraker (Creality K1, Klipper)
                try:
                    url = f"http://{host}:{port}/printer/info"
                    async with session.get(url) as resp:
                        if resp.status == 200:
                            data = await resp.json()
                            if "result" in data or "hostname" in data:
                                return PrinterType.MOONRAKER
                except:
                    pass

                # Check OctoPrint
                try:
                    url = f"http://{host}:{port}/api/version"
                    async with session.get(url) as resp:
                        if resp.status in (200, 403, 401):
                            return PrinterType.OCTOPRINT
                except:
                    pass

        except Exception as e:
            print(f"[PRINTER] Probe error for {host}:{port}: {e}")

        return PrinterType.UNKNOWN

    async def _probe_camera(self, host: str, port: int) -> Optional[str]:
        """Probe for common camera stream URLs."""
        paths = [
            "/webcam/?action=stream",
            "/webcam/stream",
            "/camera/stream",
            "/stream",
        ]

        timeout = aiohttp.ClientTimeout(total=2.0, connect=1.0)
        async with aiohttp.ClientSession(timeout=timeout) as session:
            for path in paths:
                try:
                    url = f"http://{host}:{port}{path}"
                    async with session.get(url) as resp:
                        if resp.status == 200:
                            ctype = resp.headers.get("Content-Type", "")
                            if "multipart/x-mixed-replace" in ctype or "image" in ctype:
                                print(f"[PRINTER] Found Camera: {url}")
                                return url
                except:
                    continue
        return None

    def add_printer_manually(
        self,
        name: str,
        host: str,
        port: int = 80,
        printer_type: str = "octoprint",
        api_key: Optional[str] = None,
        camera_url: Optional[str] = None,
    ) -> Printer:
        """Manually add a printer (useful when mDNS discovery fails)."""
        ptype = (
            PrinterType(printer_type)
            if printer_type in [e.value for e in PrinterType]
            else PrinterType.UNKNOWN
        )
        printer = Printer(
            name=name,
            host=host,
            port=port,
            printer_type=ptype,
            api_key=api_key,
            camera_url=camera_url,
        )
        self.printers[host] = printer
        print(f"[PRINTER] Manually added: {name} at {host}:{port}")
        return printer

    def _resolve_printer(self, target: str) -> Optional[Printer]:
        """Resolve printer by name or host."""
        if target in self.printers:
            return self.printers[target]

        for printer in self.printers.values():
            if printer.name.lower() == target.lower():
                return printer

        return None

    def _resolve_file_path(
        self, path: str, root_path: Optional[str] = None
    ) -> Optional[str]:
        """Resolve file path by checking common locations."""
        # Check if absolute path exists
        if os.path.isabs(path) and os.path.exists(path):
            return path

        common_paths = []

        # Check relative to provided root path
        if root_path:
            common_paths.append(os.path.join(root_path, path))
            if not path.startswith("cad/"):
                common_paths.append(os.path.join(root_path, "cad", path))

        # Check relative to current dir
        common_paths.append(path)
        basename = os.path.basename(path)
        common_paths.extend(
            [
                os.path.join("..", basename),
                os.path.join("cad", basename),
                "output.stl",
                "../output.stl",
            ]
        )

        for p in common_paths:
            if os.path.exists(p):
                print(f"[PRINTER] Resolved {path} -> {p}")
                return p

        return None

    async def slice_stl(
        self,
        stl_path: str,
        output_path: Optional[str] = None,
        profile_path: Optional[str] = None,
        progress_callback: Optional[Any] = None,
        root_path: Optional[str] = None,
        printer_name: Optional[str] = None,
    ) -> Optional[str]:
        """
        Slice an STL file to G-code using OrcaSlicer/PrusaSlicer CLI.

        Args:
            stl_path: Path to input STL file
            output_path: Optional output G-code path
            profile_path: Optional path to profile file (legacy)
            progress_callback: Optional async callback for progress updates
            root_path: Optional root directory to resolve relative paths
            printer_name: Optional printer name for auto-detecting profiles

        Returns:
            Path to generated G-code file, or None on failure
        """
        if not self.slicer_path:
            print("[PRINTER] Error: Slicer not found")
            return None

        # Robust path resolution
        resolved_path = self._resolve_file_path(stl_path, root_path)
        if not resolved_path:
            print(f"[PRINTER] Error: STL file not found: {stl_path}")
            return None
        stl_path = resolved_path

        # Default output path
        if not output_path:
            if root_path:
                gcode_dir = os.path.join(root_path, "gcode")
                os.makedirs(gcode_dir, exist_ok=True)
                basename = os.path.splitext(os.path.basename(stl_path))[0]
                output_path = os.path.join(gcode_dir, f"{basename}.gcode")
            else:
                output_path = stl_path.rsplit(".", 1)[0] + ".gcode"

        # Build command
        is_orca = "OrcaSlicer" in self.slicer_path

        if is_orca:
            output_dir = (
                os.path.dirname(output_path) or os.path.dirname(stl_path) or "."
            )

            cmd = [
                self.slicer_path,
                "--slice",
                "0",
                "--outputdir",
                output_dir,
            ]

            # Auto-detect profiles if printer_name is provided
            profiles = None
            if printer_name:
                profiles = self.get_profiles_for_printer(printer_name)

            settings_files = []
            if profiles:
                if profiles.get("machine"):
                    settings_files.append(profiles["machine"])
                if profiles.get("process"):
                    settings_files.append(profiles["process"])

            if settings_files:
                cmd.extend(["--load-settings", ";".join(settings_files)])

            if profiles and profiles.get("filament"):
                cmd.extend(["--load-filaments", profiles["filament"]])

            cmd.append(stl_path)

        else:
            # PrusaSlicer CLI
            cmd = [
                self.slicer_path,
                "--export-gcode",
                "--output",
                output_path,
                stl_path,
            ]

        print(f"[PRINTER] Slicing: {stl_path}")
        print(f"[PRINTER] Command: {' '.join(cmd)}")

        try:
            if progress_callback:
                await progress_callback(5, "Starting slicer...")

            result = await asyncio.to_thread(
                subprocess.run, cmd, capture_output=True, text=True
            )

            if progress_callback:
                await progress_callback(90, "Finalizing...")

            if result.returncode == 0:
                # Handle OrcaSlicer output naming
                if is_orca:
                    output_dir = os.path.dirname(output_path) or "."
                    gcode_files = glob.glob(os.path.join(output_dir, "plate_*.gcode"))

                    if not gcode_files:
                        base_name = os.path.splitext(os.path.basename(stl_path))[0]
                        expected_gcode = os.path.join(output_dir, f"{base_name}.gcode")
                        if os.path.exists(expected_gcode):
                            gcode_files = [expected_gcode]

                    if gcode_files:
                        actual_gcode = gcode_files[0]
                        if actual_gcode != output_path:
                            shutil.move(actual_gcode, output_path)

                print(f"[PRINTER] Slicing complete: {output_path}")
                if progress_callback:
                    await progress_callback(100, "Slicing Complete")
                return output_path
            else:
                print(f"[PRINTER] Slicing failed: {result.stderr}")
                return None

        except Exception as e:
            print(f"[PRINTER] Slicing error: {e}")
            return None

    async def upload_gcode(
        self, target: str, gcode_path: str, start_print: bool = False
    ) -> bool:
        """
        Upload G-code to printer and optionally start print.

        Args:
            target: Printer name or host
            gcode_path: Path to G-code file
            start_print: Whether to start printing immediately

        Returns:
            True on success, False on failure
        """
        printer = self._resolve_printer(target)
        if not printer:
            print(f"[PRINTER] Error: Printer not found: {target}")
            return False

        if not os.path.exists(gcode_path):
            print(f"[PRINTER] Error: G-code file not found: {gcode_path}")
            return False

        if printer.printer_type == PrinterType.OCTOPRINT:
            return await self._upload_octoprint(printer, gcode_path, start_print)
        elif printer.printer_type == PrinterType.MOONRAKER:
            return await self._upload_moonraker(printer, gcode_path, start_print)
        else:
            print(f"[PRINTER] Error: Unsupported printer type: {printer.printer_type}")
            return False

    async def _upload_octoprint(
        self, printer: Printer, gcode_path: str, start_print: bool
    ) -> bool:
        """Upload to OctoPrint."""
        url = f"http://{printer.host}:{printer.port}/api/files/local"
        headers = {}
        if printer.api_key:
            headers["X-Api-Key"] = printer.api_key

        filename = os.path.basename(gcode_path)

        try:
            async with aiohttp.ClientSession() as session:
                with open(gcode_path, "rb") as f:
                    data = aiohttp.FormData()
                    data.add_field("file", f, filename=filename)
                    if start_print:
                        data.add_field("print", "true")

                    async with session.post(url, data=data, headers=headers) as resp:
                        if resp.status in (200, 201, 202, 204):
                            print(f"[PRINTER] Uploaded {filename} to OctoPrint")
                            return True
                        else:
                            print(f"[PRINTER] OctoPrint upload failed ({resp.status})")
                            return False
        except Exception as e:
            print(f"[PRINTER] OctoPrint upload error: {e}")
            return False

    async def _upload_moonraker(
        self, printer: Printer, gcode_path: str, start_print: bool
    ) -> bool:
        """Upload to Moonraker."""
        url = f"http://{printer.host}:{printer.port}/server/files/upload"
        filename = os.path.basename(gcode_path)

        try:
            async with aiohttp.ClientSession() as session:
                with open(gcode_path, "rb") as f:
                    data = aiohttp.FormData()
                    data.add_field("file", f, filename=filename)

                    async with session.post(url, data=data) as resp:
                        if resp.status in (200, 201):
                            print(f"[PRINTER] Uploaded {filename} to Moonraker")

                            if start_print:
                                print_url = f"http://{printer.host}:{printer.port}/printer/print/start"
                                data_print = {"filename": filename}
                                async with session.post(
                                    print_url, json=data_print
                                ) as resp_print:
                                    if resp_print.status == 200:
                                        print(f"[PRINTER] Started print")
                                        return True
                                    else:
                                        print(
                                            f"[PRINTER] Start print failed ({resp_print.status})"
                                        )
                                        return False
                            return True
                        else:
                            print(f"[PRINTER] Moonraker upload failed ({resp.status})")
                            # Fallback to OctoPrint API
                            return await self._upload_octoprint(
                                printer, gcode_path, start_print
                            )

        except Exception as e:
            print(f"[PRINTER] Moonraker upload error: {e}")
            return False

    async def get_print_status(self, target: str) -> Optional[PrintStatus]:
        """Get current status of a printer."""
        printer = self._resolve_printer(target)
        if not printer:
            return None

        if printer.printer_type == PrinterType.OCTOPRINT:
            return await self._status_octoprint(printer)
        elif printer.printer_type == PrinterType.MOONRAKER:
            return await self._status_moonraker(printer)
        else:
            return None

    async def _status_octoprint(self, printer: Printer) -> Optional[PrintStatus]:
        """Get status from OctoPrint."""
        job_url = f"http://{printer.host}:{printer.port}/api/job"
        printer_url = f"http://{printer.host}:{printer.port}/api/printer"
        headers = {}
        if printer.api_key:
            headers["X-Api-Key"] = printer.api_key

        try:
            async with aiohttp.ClientSession() as session:
                job_data = {}
                async with session.get(job_url, headers=headers) as resp:
                    if resp.status == 200:
                        job_data = await resp.json()

                temps = {}
                async with session.get(printer_url, headers=headers) as resp:
                    if resp.status == 200:
                        printer_data = await resp.json()
                        temp_data = printer_data.get("temperature", {})
                        if "tool0" in temp_data:
                            temps["hotend"] = {
                                "current": temp_data["tool0"].get("actual", 0),
                                "target": temp_data["tool0"].get("target", 0),
                            }
                        if "bed" in temp_data:
                            temps["bed"] = {
                                "current": temp_data["bed"].get("actual", 0),
                                "target": temp_data["bed"].get("target", 0),
                            }

                if job_data:
                    progress = job_data.get("progress", {})
                    job = job_data.get("job", {})

                    return PrintStatus(
                        printer=printer.name,
                        state=job_data.get("state", "unknown").lower(),
                        progress_percent=progress.get("completion") or 0,
                        time_remaining=self._format_time(progress.get("printTimeLeft")),
                        time_elapsed=self._format_time(progress.get("printTime")),
                        filename=job.get("file", {}).get("name"),
                        temperatures=temps,
                    )

        except Exception as e:
            print(f"[PRINTER] OctoPrint status error: {e}")

        return None

    async def _status_moonraker(self, printer: Printer) -> Optional[PrintStatus]:
        """Get status from Moonraker."""
        url = f"http://{printer.host}:{printer.port}/printer/objects/query?print_stats&display_status&heater_bed&extruder"

        try:
            async with aiohttp.ClientSession() as session:
                async with session.get(url) as resp:
                    if resp.status == 200:
                        self._error_tracker.discard(printer.host)

                        data = await resp.json()
                        status = data.get("result", {}).get("status", {})
                        stats = status.get("print_stats", {})
                        display = status.get("display_status", {})
                        extruder = status.get("extruder", {})
                        bed = status.get("heater_bed", {})

                        return PrintStatus(
                            printer=printer.name,
                            state=stats.get("state", "unknown"),
                            progress_percent=(display.get("progress") or 0) * 100,
                            time_remaining=None,
                            time_elapsed=self._format_time(stats.get("print_duration")),
                            filename=stats.get("filename"),
                            temperatures={
                                "hotend": {
                                    "current": extruder.get("temperature", 0),
                                    "target": extruder.get("target", 0),
                                },
                                "bed": {
                                    "current": bed.get("temperature", 0),
                                    "target": bed.get("target", 0),
                                },
                            },
                        )

        except Exception as e:
            if printer.host not in self._error_tracker:
                print(f"[PRINTER] Moonraker status error: {e}")
                self._error_tracker.add(printer.host)

        return None

    def _format_time(self, seconds: Optional[float]) -> Optional[str]:
        if seconds is None:
            return None
        m, s = divmod(int(seconds), 60)
        h, m = divmod(m, 60)
        return f"{h:02d}:{m:02d}:{s:02d}"

    async def print_stl(
        self,
        stl_path: str,
        printer_name: str,
        profile_path: Optional[str] = None,
        root_path: Optional[str] = None,
    ) -> Dict[str, str]:
        """
        Orchestrate the full printing workflow: Slice -> Upload -> Print.

        Returns:
            Dict with 'status' and 'message'
        """
        print(f"[PRINTER] Starting print job for {stl_path} on {printer_name}")

        # 1. Resolve Printer
        printer = self._resolve_printer(printer_name)
        if not printer:
            return {
                "status": "error",
                "message": f"Printer '{printer_name}' not found.",
            }

        # 2. Slice STL
        gcode_path = await self.slice_stl(
            stl_path,
            profile_path=profile_path,
            root_path=root_path,
            printer_name=printer.name,
        )

        if not gcode_path:
            return {"status": "error", "message": "Slicing failed - check logs."}

        # 3. Upload & Start Print
        success = await self.upload_gcode(printer_name, gcode_path, start_print=True)

        if success:
            return {
                "status": "success",
                "message": f"Printing {os.path.basename(stl_path)} on {printer.name}",
            }
        else:
            return {"status": "error", "message": "Failed to upload/start print job."}


# =============================================================================
# Standalone Test
# =============================================================================
if __name__ == "__main__":

    async def main():
        agent = PrinterAgent()

        print("\n=== Testing Printer Discovery ===")
        printers = await agent.discover_printers(timeout=3)
        print(f"Found: {printers}")

        if printers:
            printer = printers[0]
            print(f"\n=== Testing Status for {printer['name']} ===")
            status = await agent.get_print_status(printer["host"])
            if status:
                print(f"Status: {status.to_dict()}")

    asyncio.run(main())
