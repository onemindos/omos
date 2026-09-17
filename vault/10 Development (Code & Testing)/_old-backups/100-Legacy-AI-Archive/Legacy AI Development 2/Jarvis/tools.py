from googlesearch import search
import spotipy
from spotipy.oauth2 import SpotifyOAuth

# Web search functions
def perform_web_search(query, num_results=3):
    try:
        results = list(search(query, num_results=num_results, advanced=True))
        return [{'title': r.title, 'description': r.description, 'url': r.url} for r in results]
    except Exception as e:
        print(f"Error performing web search: {e}")
        return []

def web_search(query):
    results = perform_web_search(query)
    if results:
        return "\n".join([f"Title: {r['title']}\nDescription: {r['description']}" for r in results])
    return f"No results found for '{query}'"

# Light control
def control_lights(state: bool):
    return f"Lights turned {'on' if state else 'off'}"

# Spotify controls
username = ''

def spotify_authenticate():
    """
    Authenticate with Spotify and return a Spotify client.
    """
    scope = (
        "user-read-playback-state user-modify-playback-state "
        "user-read-currently-playing user-read-playback-position "
        "user-library-read user-library-modify user-read-email user-read-private"
    )
    auth_manager = SpotifyOAuth(
        client_id="",
        client_secret="",
        redirect_uri="http://localhost:8888/callback",
        scope=scope,
        username=username
    )
    return spotipy.Spotify(auth_manager=auth_manager)


spotify = spotify_authenticate()

# ------------ SPOTIFY CONTROLS ------------ #

def get_active_device(spotify):
    """
    Get the active Spotify device. If no active device is found, return the first available device.
    """
    devices = spotify.devices()
    active_device = None

    for device in devices['devices']:
        if device['is_active']:
            active_device = device['id']
            break

    # If no active device, fallback to the first available device
    if not active_device and devices['devices']:
        active_device = devices['devices'][0]['id']
        print(f"No active device found. Using the available device: {devices['devices'][0]['name']}")

    return active_device

def spotify_control(action):
    """
    Control Spotify playback functions. If no active device is found, use the device issuing the command.
    """
    global spotify

    try:
        active_device = get_active_device(spotify)
        if not active_device:
            return "No Spotify devices available. Please open Spotify on a device and try again."

        if action == "get_current_playing":
            return get_current_playing_info(active_device)
        elif action == "start_music":
            spotify.start_playback(device_id=active_device)
            return "Music playback started."
        elif action == "stop_music":
            spotify.pause_playback(device_id=active_device)
            return "Music playback paused."
        elif action == "skip_to_next":
            spotify.next_track(device_id=active_device)
            return "Skipped to the next track."
        elif action == "skip_to_previous":
            spotify.previous_track(device_id=active_device)
            return "Skipped to the previous track."
        else:
            return f"Unknown Spotify action: {action}"

    except Exception as e:
        return f"Error with Spotify {action}: {str(e)}"

def get_current_playing_info(device_id):
    """
    Get the current track information from Spotify.
    """
    global spotify

    try:
        current_track = spotify.current_user_playing_track()
        if current_track is None:
            return "No track is currently playing."

        artist_name = current_track['item']['artists'][0]['name']
        album_name = current_track['item']['album']['name']
        track_title = current_track['item']['name']

        return f"Currently playing: {track_title} by {artist_name} from {album_name}."
    except spotipy.SpotifyException as e:
        return f"Error getting current track: {str(e)}"


def start_music():
    global spotify
    try:
        spotify.start_playback()
    except spotipy.SpotifyException as e:
        return f"Error in starting playback: {str(e)}"
    
def stop_music():
    global spotify
    try:
        spotify.pause_playback()
    except spotipy.SpotifyException as e:
        return f"Error in starting playback: {str(e)}"
    
def skip_to_next():
    global spotify
    try:
        spotify.next_track()
    except spotipy.SpotifyException as e:
        return f"Error in starting playback: {str(e)}"
    
def skip_to_previous():
    global spotify
    try:
        spotify.previous_track()
    except spotipy.SpotifyException as e:
        return f"Error in starting playback: {str(e)}"
    
    # Silent mode state
silent_mode_state = False

def toggle_silent_mode(state):
    global silent_mode_state
    silent_mode_state = state
    return f"Silent mode {'enabled' if state else 'disabled'}."

def is_silent_mode():
    return silent_mode_state