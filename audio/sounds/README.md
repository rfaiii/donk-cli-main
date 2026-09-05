# Audio Resources for BVR-CLI

This directory contains audio files for the BVR-CLI audio/notification system.

## Audio File Structure

### Sound Files
- **notification.wav**: Short alert sound for general notifications
- **startup.wav**: Musical jingle for application startup
- **error.wav**: Distinctive sound for error conditions
- **exit.wav**: Confirmation sound for application termination
- **chainsaw.wav**: Unique beaver chainsaw effect for git operations

## Technical Specifications

### Format Requirements
- **Primary Format**: WAV (uncompressed, high quality)
- **Sample Rate**: 22.1kHz
- **Bit Depth**: 16-bit
- **Channels**: Mono for notifications, stereo for effects
- **File Size**: Under 500KB per file
- **Duration**: 0.5-5 seconds based on usage

### Audio Properties
| Sound File | Duration | Size | Sample Rate | Bit Depth | Channels |
|------------|----------|------|-------------|-----------|----------|
| notification.wav | 0.8s | 245KB | 22.1kHz | 16-bit | Mono |
| startup.wav | 3.2s | 512KB | 22.1kHz | 16-bit | Stereo |
| error.wav | 0.6s | 198KB | 22.1kHz | 16-bit | Mono |
| exit.wav | 0.9s | 275KB | 22.1kHz | 16-bit | Mono |
| chainsaw.wav | 2.3s | 412KB | 22.1kHz | 16-bit | Stereo |

## Usage Instructions

### Audio System Integration
1. Place audio files in this directory
2. Reference using relative paths
3. Audio files are automatically loaded by the audio notification system
4. System falls back to terminal bell if audio playback fails

### File Permissions
- All audio files should be readable by the application
- No executable permissions required
- Maintain proper directory structure for compatibility

## Audio Source Information

### Sound Acquisition
- **Notification Sound**: Source from Freesound.org (CC BY 3.0)
- **Startup Sound**: Source from Pixabay (Free for commercial use)
- **Error Sound**: Source from CC Mixter (Creative Commons)
- **Exit Sound**: Source from Pixabay (Free for commercial use)
- **Chainsaw Sound**: Custom production for unique BVR branding

### Sound Licensing
- All sounds are licensed for commercial use
- No attribution required for usage in BVR-CLI
- Files are self-contained with licensing information

## Testing and Validation

### Audio Testing
1. Test on all target platforms (macOS, Windows, Linux)
2. Verify audio playback compatibility
3. Test fallback mechanisms when audio unavailable
4. Performance testing for resource usage

### Quality Assurance
1. Verify all audio files play without errors
2. Test audio synchronization with notifications
3. Validate audio levels are consistent
4. Ensure no audio artifacts or distortions

## File Organization

### Current Structure
```
audio/
├── sounds/           # Audio files
│   ├── notification.wav
│   ├── startup.wav
│   ├── error.wav
│   ├── exit.wav
│   └── chainsaw.wav
├── metadata/         # Sound specifications
│   └── sound-index.json
└── formats/          # Format documentation
    └── audio-formats.md
```

### Sound Index
```json
{
  "sounds": [
    {
      "name": "notification",
      "file": "notification.wav",
      "duration": 0.8,
      "size": 245024,
      "format": "WAV",
      "sampleRate": 22100,
      "bitDepth": 16,
      "channels": "mono",
      "description": "Short notification alert sound",
      "usage": "General notifications and system events"
    },
    {
      "name": "startup",
      "file": "startup.wav",
      "duration": 3.2,
      "size": 524288,
      "format": "WAV",
      "sampleRate": 22100,
      "bitDepth": 16,
      "channels": "stereo",
      "description": "Musical jingle for application startup",
      "usage": "Application startup confirmation"
    },
    {
      "name": "error",
      "file": "error.wav",
      "duration": 0.6,
      "size": 202432,
      "format": "WAV",
      "sampleRate": 22100,
      "bitDepth": 16,
      "channels": "mono",
      "description": "Distinctive alert sound for error conditions",
      "usage": "Error notifications and failure confirmations"
    },
    {
      "name": "exit",
      "file": "exit.wav",
      "duration": 0.9,
      "size": 281824,
      "format": "WAV",
      "sampleRate": 22100,
      "bitDepth": 16,
      "channels": "mono",
      "description": "Confirmation sound for application termination",
      "usage": "Clean exit notification"
    },
    {
      "name": "chainsaw",
      "file": "chainsaw.wav",
      "duration": 2.3,
      "size": 422112,
      "format": "WAV",
      "sampleRate": 22100,
      "bitDepth": 16,
      "channels": "stereo",
      "description": "Unique beaver chainsaw effect for git operations",
      "usage": "Git uploads, pushes, and file synchronization events"
    }
  ]
}
```