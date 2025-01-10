# 🧭 gyromidi
A server that turns your gyroscope-having mobile device into a mappable midi controller.

# Setup

# Windows
- Install the [Chocolatey Package Manager](https://chocolatey.org/)
- Install [loopMIDI](https://www.tobias-erichsen.de/software/loopmidi.html)
- Run `choco install openssl && choco install go` in an elevated powershell.
- Open loopMIDI and create a new port with the exact name "GyroMidi" (without quotes).
- Run `run_windows.bat` (you will run this any time you want to run GyroMidi).

# Unix (Mac/Linux)
- Install OpenSSL and Golang with your package manager of choice. (likely [Homebrew](https://brew.sh) for Mac)
- Run `run_unix.sh` (you will run this any time you want to run GyroMidi).

---

After running the server:
- Go to the URL printed in the console, on your gyroscopic mobile device.
- Polling mode "Polling Rate" will send gyroscope data every `x` milliseconds where `x` is the number entered in Polling Rate.
- Polling mode "On Movement" will send gyroscope data whenever the browser gets the data, (as fast as possible).
- Then turn on Polling, go to your DAW of choice, and GyroMidi will be a MIDI source.
- The Gyroscope X, Y, and Z axis will be changing the MIDI CC values set in `config.toml`.