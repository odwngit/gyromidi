# 🧭 gyromidi
A server that turns your gyroscope-having mobile device into a mappable midi controller.

# Setup
- Install OpenSSL, and the golang cli for your platform.  
- Clone this repository via `git clone https://github.com/odwngit/gyromidi && cd gyromidi` . 
- `cd ssl && ./ssl_generator.ssh` (this may look a little different depending on your platform, just run the script) . 
- `cd .. && go run .` . 
- This has now started the server on the computer that you want to use the MIDI on.  
- Now, go to the url printed out in the console on your mobile device, enter your polling rate (sensible is around ~30hz) and click start polling.  
- It is now running as a virtual midi device on your computer.  