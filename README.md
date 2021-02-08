# Home Theater Control Software

For years I used iRule software to manage my home theater. The company that produced the software was bought some years ago, but iRule continued to work... until this past Thursday, when i could no longer login. Not finding a suitable replacement, I decided to write my own software. This will probably not be of use to anyone unless they have exactly my configuration and workflow, however, I hope that it might provide a starting point to someone.

![screenshot](./screenshot.png)

## My Theater Configuration

This solution provides management for the following resources...

- __6 LimitlessLED Lights__: These are mounted in-ceiling and separated into 4 zones. This software allows me to set them to full brightness, turn them off, or set them to "dining mode". Dining mode brings the 2 lights directly overhead the primary seating positions to the lowest light setting. Whether I am dining or not, I generally keep the theater at this lighting level as it provides some ambient light without throwing light on the screen. I control the lights via a controller that accepts UDP commands and then wirelessly transmits to the lights.

- __Epson 5050UB Projector__: This is a 3-chip LCD projector that looks great. I am managing it via a Global Cache GC-100-6 controller. This controller accepts TCP commands and translates them into IR codes.

- __Denon AVR-S960H__: Frankly since I don't need to switch video channels, upscale video, or any of the other advanced features, this is maybe a bit overkill, but I love the way Denon receivers sound, especially when paired with Klipsch speakers. I control the Denon using TCP commands.

- __Roku 4630X Media Player__: Its the only media player I need, it can stream everything, even AppleTV or AirPlay. Combined with the Epson 5050, I can stream 4K HDR video. I control the Roku using HTTP REST commands.

## Design

Here are some of my thoughts on the design.

- I debated the right architecture for this solution. The iRule software was a native iPad app that handled the TCP, UDP, HTTP, etc. communication directly in the app. Ultimately I decided to build this as a web app with an API backend. The API backend will talk to the various components. This has a couple of advantages...

  - The theater room doesn't get a great WiFi signal (all the equipment is wired) and so allowing the more critical communication to happen over ethernet improves the reliability. In particular, I have noticed the lights (which are UDP) more reliably get to the desired brightness. Also, I believe the iRule software was keeping the TCP connections open and due to the unreliable network, the receiver would sometime be slow to respond.

  - Having a separate API allows for a broader range of control options, should I decide on more options in the future. For instance, it might be nice to have a physical remote that could talk to the API over a gateway.

  - It is easier to debug and update solutions that use this design.

  - Folks visiting do not need to install an app to control the theater, they just need a table. Perhaps I could write a more compact design for a phone.

  There is also one big disadvantage...

  - The API backend has to be running or the solution is of no use.

- The API backend is deployed as a docker container on my iMac. The container is set to restart=always. This should ensure it starts back up after a reboot or crash.

- I decided to focus the screen around my common workflows -> change to a channel on Roku, browse or search content, play the video, maybe pause/rewind, go back to home.

  - In iRule, the screens I designed around each device (Roku, AppleTV, Bluray player, etc.). I also had built out full functionality replacing all remotes, but in reality it is a very rare case that any of that was used.

  - I have an AppleTV and a Bluray player also connected to the receiver, but I haven't used them in years, so I didn't bother to add them to the controller.

## Usage

Here are some things I can do with the new controller software...

- __System Power On__: This sets the lighting to dining mode, turns on the receiver, and turns on the projector. It also queries to find out the current volume setting of the receiver.

- __System Standby__: This sets the lighting to full brightness, sets the receiver to standby, and sets the projector to standby.

- __Touchpad__: The main section of the screen is a touchpad. I can swipe my finger up, down, left, or right to navigate or press to select something.

- __Roku Controls__: This surface should allow you to do anything you can do with a Roku remote, such as return to home, go back, get info, etc.

- __Volume__: Controlling the volume or muting is available at the lower right. It also displays the current volume, though changing the volume will show that on the main screen anyway.

- __Common Channels__: The common channels I watch are displayed on the left side. Clicking on any of them will take me into that channel on the Roku.

- __Text Input__: If I need to search for something, type a password, or whatever else, I can type that in and then click the "enter" button to send it to the Roku. This also serves as a command line for functions that don't have a button. Currently, you can change inputs by typing something like `$input:roku`.

- __Overrides__: Sometimes I might want to override the projector, receiver, or lighting to turn them on or off or dim. There are buttons to handle these more frequent uncommon operations.

## Software Configuration

For the software to work, some environmental variables must be set. These can be configured in a .env file, using ENV in the Dockerfile, or any other method.

- __PORT__: [DEFAULT: 8080] This determines the port the API backend is running on.

- __PROJECTOR__: [REQUIRED] This should be the IPaddress:port or FQDN:port of the Global Cache controller connected to the projector via IR.

- __PROJECTOR_IR_PORT__: [REQUIRED] This is the IR port number on the Global Cache controller that is plugged into the projector.

- __RECEIVER__: [REQUIRED] This should be the IPaddress:port or FQDN:port of the Denon receiver. Denon receivers use port 23 for TCP commands.

- __ROKU__: [REQUIRED] This should be the IPaddress:port or FQDN:port of the Roku. Rokus use port 8060 for HTTP REST commands.

- __LIGHTS__: [REQUIRED] This should be the IPaddress:port or FQDN:port of the lighting controller. LimitlessLED uses port 8899 for UDP commands.

## Pictures

Here is what the room looks like with all lights on. There are 2 rows of 3 lights each that run the length of the room.

![theater](./theater.png)

Here is the equipment closet.

![equipment](./equipment.png)

While taking a picture of a projected image doesn't do it justice, here it is anyway.

![screen in use](./screen.png)

I wanted to take a picture of the room with the dining room lights on, but without using a flash it was really dark and grainy. In that mode, the 2 center lights are on at the lowest setting (10%), it provides a nice subtle light when seated it the main positions without filling the room with too much ambient light.
