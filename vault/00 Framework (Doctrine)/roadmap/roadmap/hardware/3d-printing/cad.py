"""
CAD Agent - 3D Model Generation with build123d
Extracted from ADA V2 (https://github.com/nazirlouis/ada_v2)

This agent generates 3D models from natural language descriptions using:
- Gemini API for code generation (adapt to Claude/Bedrock for OneMind)
- build123d library for parametric CAD
- Local Python execution for model generation
- STL export for downstream use

Dependencies:
    pip install build123d cadquery-ocp google-genai python-dotenv

For OneMind Integration:
    - Replace Gemini with Claude Opus via AWS Bedrock
    - Use NATS for status callbacks instead of direct callbacks
    - Integrate with hardware/ directory for robotics projects
"""

import os
import json
import asyncio
import re
import base64
import subprocess
import sys
import tempfile
from datetime import datetime
from typing import Optional, Callable, Dict, Any

# For OneMind: Replace with your LLM client
# from google import genai
# from google.genai import types

from dotenv import load_dotenv

load_dotenv()


class CadAgent:
    """
    3D CAD Generation Agent using build123d library.

    Generates Python scripts that create 3D models, executes them locally,
    and returns the resulting STL file.
    """

    # System instruction for the LLM
    SYSTEM_INSTRUCTION = """
You are a Python-based 3D CAD Engineer using the `build123d` library.
Your goal is to write a Python script that generates a 3D model based on the user's request.

Requirements:
1. Start with `from build123d import *`.
2. Include `import numpy as np` if you use any numpy functions (like `np.sign`, `np.pi`).
3. You MUST assign the final object to a variable named `result_part`.
4. If you create a sketch or line, extrude it to make it a solid `Part`.
5. The model should be centered at (0,0,0) and have reasonable dimensions (mm).
6. **IMPORTANT**: Do NOT use old or PascalCase function names for core operations.
   - Use `make_face()` instead of `MakeFace()`.
   - Use `extrude()` instead of `Extrude()`.
   - Use `fillet()` instead of `Fillet()`.
   - Use `chamfer()` instead of `Chamfer()`.
   - Use `revolve()` instead of `Revolve()`.
   - Use `loft()` instead of `Loft()`.
   - Use `sweep()` instead of `Sweep()`.
   - Use `offset()` instead of `Offset()`.
   - generally prefer lowercase builder methods inside contexts.

7. **Vector Access**: Do NOT access vector components like `v.X`, `v.Y`, `v.Z` unless you are sure they exist (use `v.X` etc on Vector objects, but ensure they are Vectors).
8. **Final Output**: The script MUST end by exporting the final part to an STL file named 'output.stl'.
   - `export_stl(result_part, 'output.stl')`

9. **Robustness**: Operations like `fillet()` and `chamfer()` will crash if the radius is too large for the geometry.
   - Use conservative values (e.g., 0.5mm to 2mm) unless you are certain of the dimensions.
   - If a fillet is purely aesthetic, keep it small to ensure success.

Example Script:
```python
from build123d import *

with BuildPart() as p:
    Box(10, 10, 10)
    Fillet(p.edges(), radius=1)

result_part = p.part
export_stl(result_part, 'output.stl')
```
"""

    def __init__(
        self,
        on_thought: Optional[Callable[[str], None]] = None,
        on_status: Optional[Callable[[Dict[str, Any]], None]] = None,
        llm_client: Optional[Any] = None,
    ):
        """
        Initialize the CAD Agent.

        Args:
            on_thought: Callback for streaming LLM thoughts (for UI feedback)
            on_status: Callback for status updates (generating, retrying, failed)
            llm_client: LLM client for code generation. If None, will attempt to
                       use Gemini. For OneMind, pass your Bedrock/Claude client.
        """
        self.on_thought = on_thought
        self.on_status = on_status
        self.llm_client = llm_client

        # For OneMind: Replace with Claude model
        self.model = "gemini-3-pro-preview"

        # Initialize Gemini client if no custom client provided
        if self.llm_client is None:
            try:
                from google import genai

                self.client = genai.Client(
                    http_options={"api_version": "v1beta"},
                    api_key=os.getenv("GEMINI_API_KEY"),
                )
            except ImportError:
                print(
                    "[CadAgent] Warning: google-genai not installed. Use custom llm_client."
                )
                self.client = None

    async def _generate_code_with_gemini(self, prompt: str) -> str:
        """
        Generate build123d code using Gemini API.

        For OneMind: Replace this method with Claude/Bedrock code generation.
        """
        from google.genai import types

        raw_content = ""
        stream = await self.client.aio.models.generate_content_stream(
            model=self.model,
            contents=prompt,
            config=types.GenerateContentConfig(
                system_instruction=self.SYSTEM_INSTRUCTION,
                temperature=1.0,
                thinking_config=types.ThinkingConfig(include_thoughts=True),
            ),
        )

        async for chunk in stream:
            if (
                chunk.candidates
                and chunk.candidates[0].content
                and chunk.candidates[0].content.parts
            ):
                for part in chunk.candidates[0].content.parts:
                    if not part.text:
                        continue
                    elif part.thought:
                        # Stream thought to callback
                        if self.on_thought:
                            self.on_thought(part.text)
                    else:
                        # Accumulate answer text
                        raw_content += part.text

        return raw_content

    async def _generate_code_with_custom_llm(self, prompt: str) -> str:
        """
        Generate build123d code using custom LLM client.

        Override this method for OneMind's Claude/Bedrock integration.
        """
        if self.llm_client is None:
            raise ValueError("No LLM client configured")

        # Example for OneMind's Claude integration:
        # response = await self.llm_client.messages.create(
        #     model="claude-opus-4-5-20250227",
        #     max_tokens=4096,
        #     system=self.SYSTEM_INSTRUCTION,
        #     messages=[{"role": "user", "content": prompt}]
        # )
        # return response.content[0].text

        raise NotImplementedError("Implement custom LLM code generation for your setup")

    def _extract_python_code(self, raw_content: str) -> Optional[str]:
        """Extract Python code block from LLM response."""
        if not raw_content:
            print("[CadAgent DEBUG] [ERR] Empty response from model.")
            return None

        # Try to extract ```python block
        code_match = re.search(r"```python(.*?)```", raw_content, re.DOTALL)
        if code_match:
            return code_match.group(1).strip()

        # Fallback: assume entire text is code if it looks like build123d code
        print("[CadAgent DEBUG] [WARN] No ```python block found. Trying heuristic...")
        if "build123d" in raw_content or "BuildPart" in raw_content:
            return raw_content

        print("[CadAgent DEBUG] [ERR] Could not extract python code.")
        return None

    async def _execute_script(self, script_path: str) -> tuple:
        """
        Execute the generated Python script.

        Returns:
            Tuple of (success: bool, stdout: str, stderr: str)
        """
        try:
            proc = await asyncio.to_thread(
                subprocess.run,
                [sys.executable, script_path],
                capture_output=True,
                text=True,
            )
            return proc.returncode == 0, proc.stdout, proc.stderr
        except Exception as e:
            print(f"[CadAgent DEBUG] [ERR] Subprocess run failed: {e}")
            return False, "", str(e)

    async def generate_prototype(
        self, prompt: str, output_dir: Optional[str] = None
    ) -> Optional[Dict[str, Any]]:
        """
        Generates 3D geometry by asking LLM for a script, then running it locally.

        Args:
            prompt: User's description of the model to generate.
            output_dir: Directory to save the script and STL. If None, uses temp dir.

        Returns:
            Dict with 'format', 'data' (base64 STL), and 'file_path', or None on failure.
        """
        print(f"[CadAgent DEBUG] [START] Generation started for: '{prompt}'")

        try:
            # Use provided output_dir or fall back to temp
            if output_dir:
                os.makedirs(output_dir, exist_ok=True)
                work_dir = output_dir
            else:
                work_dir = tempfile.gettempdir()

            # Generate timestamped filename
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            output_stl = os.path.join(work_dir, f"output_{timestamp}.stl")
            script_path = os.path.join(work_dir, "current_design.py")

            max_retries = 3
            current_prompt = f"You are a build123d expert. Write a generic python script to create a 3D model of: {prompt}. Ensure you export to 'output.stl'. Unscaled."

            for attempt in range(max_retries):
                print(f"[CadAgent DEBUG] Attempt {attempt + 1}/{max_retries}")

                # Emit status update
                if self.on_status:
                    status_info = {
                        "status": "generating" if attempt == 0 else "retrying",
                        "attempt": attempt + 1,
                        "max_attempts": max_retries,
                        "error": None,
                    }
                    self.on_status(status_info)

                # 1. Generate code with LLM
                if self.client:
                    raw_content = await self._generate_code_with_gemini(current_prompt)
                else:
                    raw_content = await self._generate_code_with_custom_llm(
                        current_prompt
                    )

                # 2. Extract code block
                code = self._extract_python_code(raw_content)
                if not code:
                    continue

                # 3. Save to local file
                # Fix for Windows paths: escape backslashes
                safe_output_path = output_stl.replace("\\", "\\\\")

                with open(script_path, "w") as f:
                    # Inject output path into the script
                    code_with_path = code.replace("output.stl", safe_output_path)
                    f.write(code_with_path)

                print(f"[CadAgent DEBUG] [EXEC] Running local script: {script_path}")

                # 4. Execute locally
                success, stdout, stderr = await self._execute_script(script_path)

                if not success:
                    error_msg = stderr
                    error_lines = error_msg.strip().split("\n")
                    short_error = (
                        error_lines[-1][:100] if error_lines else "Unknown error"
                    )
                    print(
                        f"[CadAgent DEBUG] [ERR] Script Execution Failed:\n{error_msg}"
                    )

                    # Emit retry status with error
                    if self.on_status:
                        self.on_status(
                            {
                                "status": "retrying",
                                "attempt": attempt + 1,
                                "max_attempts": max_retries,
                                "error": short_error,
                            }
                        )

                    # Prepare feedback for next attempt
                    current_prompt = f"""
The Python script you generated failed to execute with the following error:
{error_msg}

Please fix the code to resolve this error. Return the full corrected script. 
Ensure you still export to 'output.stl'.
Original request: {prompt}
"""
                    continue  # Retry loop

                print(f"[CadAgent DEBUG] [OK] Script executed successfully.")

                # 5. Read Output
                if os.path.exists(output_stl):
                    print(f"[CadAgent DEBUG] [file] '{output_stl}' found.")
                    with open(output_stl, "rb") as f:
                        stl_data = f.read()

                    b64_stl = base64.b64encode(stl_data).decode("utf-8")

                    return {"format": "stl", "data": b64_stl, "file_path": output_stl}
                else:
                    print(f"[CadAgent DEBUG] [ERR] '{output_stl}' was not generated.")
                    current_prompt = f"The script executed successfully but 'output.stl' was not found. Ensure you call `export_stl(result_part, 'output.stl')` at the end."
                    continue

            # If loop finishes without success
            print("[CadAgent DEBUG] [ERR] All attempts failed.")
            if self.on_status:
                self.on_status(
                    {
                        "status": "failed",
                        "attempt": max_retries,
                        "max_attempts": max_retries,
                        "error": "All generation attempts failed",
                    }
                )
            return None

        except Exception as e:
            print(f"CadAgent Error: {e}")
            import traceback

            traceback.print_exc()
            return None

    async def iterate_prototype(
        self, prompt: str, output_dir: Optional[str] = None
    ) -> Optional[Dict[str, Any]]:
        """
        Iterates on the existing design by reading 'current_design.py' and applying changes.

        Args:
            prompt: User's description of the changes to make.
            output_dir: Directory containing existing script and where to save new STL.

        Returns:
            Dict with 'format', 'data' (base64 STL), and 'file_path', or None on failure.
        """
        print(f"[CadAgent DEBUG] [START] Iteration started for: '{prompt}'")

        # Use provided output_dir or fall back to temp
        if output_dir:
            os.makedirs(output_dir, exist_ok=True)
            work_dir = output_dir
        else:
            work_dir = tempfile.gettempdir()

        # Generate timestamped filename for the output
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        script_path = os.path.join(work_dir, "current_design.py")
        output_stl = os.path.join(work_dir, f"output_{timestamp}.stl")

        existing_code = ""

        if os.path.exists(script_path):
            with open(script_path, "r") as f:
                existing_code = f.read()

            # Sanitize existing code: replace any absolute paths with 'output.stl'
            existing_code = re.sub(
                r"['\"]C:\\\\?Users\\\\?[^'\"]+\\\\?output[^'\"]*\.stl['\"]",
                "'output.stl'",
                existing_code,
            )
            # Also handle forward-slash variants
            existing_code = re.sub(
                r"['\"]C:/Users/[^'\"]+/output[^'\"]*\.stl['\"]",
                "'output.stl'",
                existing_code,
            )
        else:
            print(
                "[CadAgent DEBUG] [WARN] No existing script found. Falling back to fresh generation."
            )
            return await self.generate_prototype(prompt, output_dir)

        try:
            max_retries = 3
            current_prompt = f"""
You are iterating on an existing 3D model script.

Current Python Code:
```python
{existing_code}
```

User Request: {prompt}

Task: Rewrite the code to satisfy the user's request while maintaining the rest of the model structure.
Ensure you still export to 'output.stl'.
"""

            for attempt in range(max_retries):
                print(f"[CadAgent DEBUG] Iteration Attempt {attempt + 1}/{max_retries}")

                # Emit status update
                if self.on_status:
                    status_info = {
                        "status": "generating" if attempt == 0 else "retrying",
                        "attempt": attempt + 1,
                        "max_attempts": max_retries,
                        "error": None,
                    }
                    self.on_status(status_info)

                # 1. Generate code with LLM
                if self.client:
                    raw_content = await self._generate_code_with_gemini(current_prompt)
                else:
                    raw_content = await self._generate_code_with_custom_llm(
                        current_prompt
                    )

                # 2. Extract code block
                code = self._extract_python_code(raw_content)
                if not code:
                    continue

                # 3. Save to local file (overwrite so next iteration builds on this)
                safe_output_path = output_stl.replace("\\", "\\\\")

                with open(script_path, "w") as f:
                    code_with_path = code.replace("output.stl", safe_output_path)
                    f.write(code_with_path)

                print(f"[CadAgent DEBUG] [EXEC] Running local script: {script_path}")

                # 4. Execute locally
                success, stdout, stderr = await self._execute_script(script_path)

                if not success:
                    error_msg = stderr
                    print(
                        f"[CadAgent DEBUG] [ERR] Script Execution Failed:\n{error_msg}"
                    )

                    current_prompt = f"""
The updated Python script you generated failed to execute with the following error:
{error_msg}

Please fix the code to resolve this error. Return the full corrected script. 
Ensure you still export to 'output.stl'.
"""
                    continue

                print(f"[CadAgent DEBUG] [OK] Script executed successfully.")

                # 5. Read Output
                if os.path.exists(output_stl):
                    print(f"[CadAgent DEBUG] [file] '{output_stl}' found.")
                    with open(output_stl, "rb") as f:
                        stl_data = f.read()

                    b64_stl = base64.b64encode(stl_data).decode("utf-8")

                    return {"format": "stl", "data": b64_stl, "file_path": output_stl}
                else:
                    print(f"[CadAgent DEBUG] [ERR] '{output_stl}' was not generated.")
                    current_prompt = f"The script executed successfully but '{output_stl}' was not found. Ensure you call `export_stl(result_part, 'output.stl')` at the end."
                    continue

            # If loop finishes without success
            print("[CadAgent DEBUG] [ERR] All attempts failed.")
            if self.on_status:
                self.on_status(
                    {
                        "status": "failed",
                        "attempt": max_retries,
                        "max_attempts": max_retries,
                        "error": "All iteration attempts failed",
                    }
                )
            return None

        except Exception as e:
            print(f"CadAgent Error: {e}")
            import traceback

            traceback.print_exc()
            return None


# =============================================================================
# Standalone Test
# =============================================================================
if __name__ == "__main__":

    async def main():
        agent = CadAgent()

        print("\n=== Testing CAD Generation ===")
        result = await agent.generate_prototype(
            "A simple gear with 12 teeth, 20mm outer diameter, 5mm thick"
        )

        if result:
            print(f"Generated STL: {result['file_path']}")
            print(f"Base64 size: {len(result['data'])} chars")
        else:
            print("Generation failed")

    asyncio.run(main())
