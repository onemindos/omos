"""
Face Authenticator - Local Biometric Authentication via MediaPipe
Extracted from ADA V2 (https://github.com/nazirlouis/ada_v2)

This agent provides face-based authentication using:
- MediaPipe Face Landmarker for face detection
- Cosine similarity for face matching
- Local processing with no cloud dependency

Dependencies:
    pip install mediapipe opencv-python numpy

For OneMind Integration:
    - Use for local identity verification before sensitive operations
    - Can run alongside your existing auth systems
    - Good for HITL (Human-in-the-Loop) confirmations
"""

import os
import base64
import asyncio
from typing import Optional, Callable
import urllib.request

import cv2
import numpy as np

# MediaPipe imports
import mediapipe as mp
from mediapipe.tasks import python as mp_python
from mediapipe.tasks.python import vision


class FaceAuthenticator:
    """
    Face authentication using MediaPipe Face Landmarker.

    Compares live camera feed against a reference face image
    using landmark-based cosine similarity.
    """

    # MediaPipe Face Landmarker model URL
    MODEL_URL = "https://storage.googleapis.com/mediapipe-models/face_landmarker/face_landmarker/float16/1/face_landmarker.task"
    MODEL_PATH = os.path.join(os.path.dirname(__file__), "face_landmarker.task")

    def __init__(
        self,
        reference_image_path: str = "reference.jpg",
        on_status_change: Optional[Callable] = None,
        on_frame: Optional[Callable] = None,
        similarity_threshold: float = 0.15,
    ):
        """
        Initialize the Face Authenticator.

        Args:
            reference_image_path: Path to the user's reference photo.
            on_status_change: Async callback(is_authenticated: bool) for auth state changes.
            on_frame: Async callback(frame_data_b64: str) to send frames to frontend.
            similarity_threshold: Cosine similarity threshold (lower = stricter matching).
        """
        self.reference_image_path = reference_image_path
        self.on_status_change = on_status_change
        self.on_frame = on_frame
        self.similarity_threshold = similarity_threshold

        self.authenticated = False
        self.running = False
        self.reference_landmarks = None
        self.landmarker = None

        self._ensure_model()
        self._init_landmarker()
        self._load_reference()

    def _ensure_model(self):
        """Download the MediaPipe Face Landmarker model if not present."""
        if not os.path.exists(self.MODEL_PATH):
            print(f"[AUTH] Downloading Face Landmarker model...")
            try:
                urllib.request.urlretrieve(self.MODEL_URL, self.MODEL_PATH)
                print(f"[AUTH] [OK] Model downloaded to {self.MODEL_PATH}")
            except Exception as e:
                print(f"[AUTH] [ERR] Failed to download model: {e}")

    def _init_landmarker(self):
        """Initialize the MediaPipe Face Landmarker."""
        if not os.path.exists(self.MODEL_PATH):
            print("[AUTH] [ERR] Face Landmarker model not found. Cannot initialize.")
            return

        try:
            base_options = mp_python.BaseOptions(model_asset_path=self.MODEL_PATH)
            options = vision.FaceLandmarkerOptions(
                base_options=base_options,
                output_face_blendshapes=False,
                output_facial_transformation_matrixes=False,
                num_faces=1,
            )
            self.landmarker = vision.FaceLandmarker.create_from_options(options)
            print("[AUTH] [OK] Face Landmarker initialized.")
        except Exception as e:
            print(f"[AUTH] [ERR] Failed to initialize Face Landmarker: {e}")

    def _extract_landmarks(self, image_rgb: np.ndarray) -> Optional[np.ndarray]:
        """
        Extract normalized face landmarks from an RGB image.

        Returns:
            Flattened numpy array of (x, y, z) coordinates, or None if no face found.
        """
        if self.landmarker is None:
            return None

        try:
            mp_image = mp.Image(image_format=mp.ImageFormat.SRGB, data=image_rgb)
            result = self.landmarker.detect(mp_image)

            if result.face_landmarks and len(result.face_landmarks) > 0:
                landmarks = result.face_landmarks[0]
                # Convert to numpy array of (x, y, z) coordinates
                coords = np.array(
                    [[lm.x, lm.y, lm.z] for lm in landmarks], dtype=np.float32
                )
                return coords.flatten()
            return None
        except Exception as e:
            print(f"[AUTH] [ERR] Landmark extraction failed: {e}")
            return None

    def _compare_landmarks(
        self, landmarks1: Optional[np.ndarray], landmarks2: Optional[np.ndarray]
    ) -> bool:
        """
        Compare two landmark vectors using cosine similarity.

        Returns:
            True if similarity is above threshold (1 - similarity_threshold).
        """
        if landmarks1 is None or landmarks2 is None:
            return False

        # Normalize vectors
        norm1 = np.linalg.norm(landmarks1)
        norm2 = np.linalg.norm(landmarks2)

        if norm1 == 0 or norm2 == 0:
            return False

        # Cosine similarity
        similarity = np.dot(landmarks1, landmarks2) / (norm1 * norm2)

        # Threshold check (similarity should be close to 1 for a match)
        is_match = similarity > (1 - self.similarity_threshold)
        if is_match:
            print(f"[AUTH] Face match! Similarity: {similarity:.4f}")
        return is_match

    def _load_reference(self):
        """Load and process the reference face image."""
        if not os.path.exists(self.reference_image_path):
            print(
                f"[AUTH] [WARN] Reference file not found at {self.reference_image_path}. Authentication will fail."
            )
            return

        try:
            print("[AUTH] Loading reference image...")
            img_bgr = cv2.imread(self.reference_image_path)
            if img_bgr is None:
                print(
                    f"[AUTH] [ERR] Failed to read image file: {self.reference_image_path}"
                )
                return

            # Convert to RGB
            image_rgb = cv2.cvtColor(img_bgr, cv2.COLOR_BGR2RGB)

            self.reference_landmarks = self._extract_landmarks(image_rgb)

            if self.reference_landmarks is not None:
                print("[AUTH] [OK] Reference face landmarks extracted successfully.")
            else:
                print("[AUTH] [ERR] No face found in reference image.")
        except Exception as e:
            print(f"[AUTH] [ERR] Error loading reference: {e}")

    def set_reference_from_base64(self, image_b64: str) -> bool:
        """
        Set reference face from a base64-encoded image.

        Args:
            image_b64: Base64 encoded image data.

        Returns:
            True if reference was successfully set, False otherwise.
        """
        try:
            # Decode base64 to image
            image_data = base64.b64decode(image_b64)
            nparr = np.frombuffer(image_data, np.uint8)
            img_bgr = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

            if img_bgr is None:
                print("[AUTH] [ERR] Failed to decode base64 image")
                return False

            # Convert to RGB and extract landmarks
            image_rgb = cv2.cvtColor(img_bgr, cv2.COLOR_BGR2RGB)
            self.reference_landmarks = self._extract_landmarks(image_rgb)

            if self.reference_landmarks is not None:
                # Optionally save the reference image
                cv2.imwrite(self.reference_image_path, img_bgr)
                print("[AUTH] [OK] Reference face set from base64 image.")
                return True
            else:
                print("[AUTH] [ERR] No face found in provided image.")
                return False

        except Exception as e:
            print(f"[AUTH] [ERR] Error setting reference from base64: {e}")
            return False

    async def start_authentication_loop(self):
        """
        Start the authentication loop using webcam.

        Continuously captures frames and compares against reference
        until authentication succeeds or stop() is called.
        """
        if self.authenticated:
            print("[AUTH] Already authenticated.")
            if self.on_status_change:
                await self.on_status_change(True)
            return

        if self.reference_landmarks is None:
            print("[AUTH] [ERR] Cannot start auth loop: No reference landmarks.")
            return

        self.running = True
        print("[AUTH] Starting camera for authentication...")

        # Capture the current (main) event loop
        loop = asyncio.get_running_loop()

        # Use a separate thread for blocking camera/CV operations
        await asyncio.to_thread(self._run_cv_loop, loop)

        print("[AUTH] Authentication loop finished.")

    def stop(self):
        """Stop the authentication loop."""
        print("[AUTH] Stopping authentication loop...")
        self.running = False

    def reset(self):
        """Reset authentication state (for re-authentication)."""
        self.authenticated = False
        self.running = False

    def _run_cv_loop(self, loop):
        """
        Internal method that runs the camera loop in a separate thread.

        This is blocking and runs in asyncio.to_thread().
        """

        def try_open_camera(index):
            print(f"[AUTH] Trying to open camera with index {index}...")

            # Use appropriate backend based on platform
            import platform

            if platform.system() == "Darwin":  # macOS
                cap = cv2.VideoCapture(index, cv2.CAP_AVFOUNDATION)
            else:
                cap = cv2.VideoCapture(index)

            if not cap.isOpened():
                print(f"[AUTH] [ERR] Could not open video device {index}.")
                return None

            ret, frame = cap.read()
            if not ret:
                print(
                    f"[AUTH] [ERR] Opened device {index} but failed to read first frame."
                )
                cap.release()
                return None

            print(f"[AUTH] [OK] Successfully opened and read from device {index}.")
            return cap

        # Try multiple camera indices
        video_capture = try_open_camera(0)
        if video_capture is None:
            print("[AUTH] Device 0 failed. Trying device 1...")
            video_capture = try_open_camera(1)

        if video_capture is None:
            print(
                "[AUTH] [ERR] All camera attempts failed. Authentication cannot proceed."
            )
            self.running = False
            return

        process_this_frame = True
        frame_count = 0

        while self.running and not self.authenticated:
            ret, frame = video_capture.read()
            if not ret:
                print("[AUTH] [ERR] Failed to read frame from camera loop.")
                break

            # Convert BGR to RGB
            rgb_frame = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)

            # Process every other frame for performance
            if process_this_frame:
                current_landmarks = self._extract_landmarks(rgb_frame)

                if self._compare_landmarks(self.reference_landmarks, current_landmarks):
                    self.authenticated = True
                    print("[AUTH] [OPEN] FACE RECOGNIZED! Access Granted.")
                    if self.on_status_change:
                        asyncio.run_coroutine_threadsafe(
                            self.on_status_change(True), loop
                        )
                    self.running = False
                    break

            process_this_frame = not process_this_frame

            # Send frame to frontend if callback exists (every few frames for performance)
            if self.on_frame and frame_count % 3 == 0:
                small_frame = cv2.resize(frame, (0, 0), fx=0.5, fy=0.5)
                _, buffer = cv2.imencode(
                    ".jpg", small_frame, [cv2.IMWRITE_JPEG_QUALITY, 70]
                )
                b64_str = base64.b64encode(buffer).decode("utf-8")

                asyncio.run_coroutine_threadsafe(self.on_frame(b64_str), loop)

            frame_count += 1

        video_capture.release()

    async def authenticate_single_frame(self, image_b64: str) -> bool:
        """
        Authenticate a single frame (for non-webcam scenarios).

        Args:
            image_b64: Base64 encoded image to authenticate.

        Returns:
            True if face matches reference, False otherwise.
        """
        if self.reference_landmarks is None:
            print("[AUTH] [ERR] No reference landmarks loaded.")
            return False

        try:
            # Decode base64 to image
            image_data = base64.b64decode(image_b64)
            nparr = np.frombuffer(image_data, np.uint8)
            img_bgr = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

            if img_bgr is None:
                return False

            # Convert to RGB and extract landmarks
            image_rgb = cv2.cvtColor(img_bgr, cv2.COLOR_BGR2RGB)
            current_landmarks = self._extract_landmarks(image_rgb)

            if self._compare_landmarks(self.reference_landmarks, current_landmarks):
                self.authenticated = True
                if self.on_status_change:
                    await self.on_status_change(True)
                return True

            return False

        except Exception as e:
            print(f"[AUTH] [ERR] Single frame auth error: {e}")
            return False


# =============================================================================
# Utility: Capture Reference Face
# =============================================================================
def capture_reference_face(output_path: str = "reference.jpg"):
    """
    Opens the webcam and captures a frame when the user presses 's' or 'Space'.
    Saves the frame to the specified output path.

    This is a standalone utility for enrolling a new face.
    """
    cap = cv2.VideoCapture(0)

    if not cap.isOpened():
        print("Error: Could not open webcam.")
        return False

    print("Press 's' or 'SPACE' to capture your face.")
    print("Press 'q' or 'ESC' to quit without saving.")

    captured = False

    while True:
        ret, frame = cap.read()
        if not ret:
            print("Error: Failed to capture frame.")
            break

        # Display the resulting frame
        cv2.imshow("Capture Reference Face", frame)

        key = cv2.waitKey(1) & 0xFF

        # 's' or SPACE to save
        if key == ord("s") or key == 32:
            cv2.imwrite(output_path, frame)
            print(f"Reference face saved to {output_path}")
            captured = True
            break

        # 'q' or ESC to quit
        if key == ord("q") or key == 27:
            print("Capture cancelled.")
            break

    cap.release()
    cv2.destroyAllWindows()
    return captured


# =============================================================================
# Standalone Test
# =============================================================================
if __name__ == "__main__":

    async def main():
        print("\n=== Face Authenticator Test ===")
        print("1. Capture reference face")
        print("2. Test authentication")
        choice = input("Enter choice (1 or 2): ").strip()

        if choice == "1":
            script_dir = os.path.dirname(os.path.abspath(__file__))
            save_path = os.path.join(script_dir, "reference.jpg")
            capture_reference_face(save_path)

        elif choice == "2":
            script_dir = os.path.dirname(os.path.abspath(__file__))
            ref_path = os.path.join(script_dir, "reference.jpg")

            if not os.path.exists(ref_path):
                print(f"Reference image not found at {ref_path}")
                print("Please capture a reference face first (option 1)")
                return

            async def on_status_change(authenticated: bool):
                print(f"[CALLBACK] Authentication status: {authenticated}")

            authenticator = FaceAuthenticator(
                reference_image_path=ref_path, on_status_change=on_status_change
            )

            print("\nStarting authentication... Look at the camera.")
            await authenticator.start_authentication_loop()

            if authenticator.authenticated:
                print("\n✓ Authentication successful!")
            else:
                print("\n✗ Authentication failed or cancelled.")

        else:
            print("Invalid choice.")

    asyncio.run(main())
