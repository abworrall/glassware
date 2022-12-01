# Glassware

A project that lets you setup some physical objects, and trigger
software events when they change weight - e.g. play a particular
playlist when a stopped is removed from a bottle.

This project has a few moving pieces:
- an ardunio program, which just polls the analog pins and prints values via serial output
- a controller program that runs on a pi, reading the tty from arduino board(s), looking for significant changes
- music streaming and audio output integration

## Hardware: bill of materials

This was developed to run on the following hardware:
- raspberry pi (any Pi with networking and USB will do), and a power supply
- arduino board ([link](https://www.adafruit.com/product/2488))
   - each board can handle 6 analog FSR sensors, can have multiple boards
- force sensitive resistors (FSRs), one per bottle ([link](https://www.sparkfun.com/products/9376))
   - it's easier to establish strong contact on the bigger square sensor, and so get steadier values
- a clincher connector for each FSR ([link](https://www.sparkfun.com/products/14194))
- wiring breadboard ([link](https://www.sparkfun.com/products/12002))
- mounting board for everything ([link](https://www.sparkfun.com/products/11235))
- some 10Kohm resistors ([link](https://www.adafruit.com/product/2784))
- some wires ([short](https://www.adafruit.com/product/1956) and [long](https://www.adafruit.com/product/1955))
- a USB micro-B cable to link each ardino board back to the pi ([link](https://www.sparkfun.com/products/13244))

# Setup

## Dev machine

You'll want to setup some kind of dev machine (linux, mac, whatever):
- configure SSH client stuff (so you can SSH/SCP to the pi)
- install the [arduino IDE](https://www.arduino.cc/en/software)
- install the Go programming language, [golang](https://go.dev/)

## Get and build the glassware software

```
# Check out the code, where golang toolchain expects to see it
mkdir -p ~/go/src/github.com/abworrall
cd  ~/go/src/github.com/abworrall
git clone https://github.com/abworrall/glassware
cd glassware

# Compile the command, so you can run it locally
go build ./cmd/gw
```

To cross-compile the command to run on a Raspberry Pi, set a few ENV vars:
```
GOOS=linux GOARCH=arm go build ./cmd/gw
```

## Arduino: board setup

- Plug the arduino board into your dev machine, using the USB micro-B cable.
- Start the arduino IDE, and open up the arduino sketch file `./arduino/readpins/readpins.ino`.
- In the IDE, set [tools>board] as `Arduino UNO`, and [tools>port] as `/dev/ttyUSBS0`.
   - If you're using a different flavor of Arduino, or not using Linux, you'll need to figure this bit out yourself.
- If you're using multiple arduino boards, *make sure* to edit the `readpins` sketch each time, so that each board gets a unique controller ID.
- Click 'upload', to load the sketch into the board.
   - If this works, you'll see the red LED on the board start flashing twice a second.

You should have `gw` compiled for your dev machine, as above. You can
now run it in verbose mode, to see if it can talk to the arduino
board, and that the board is sending back the right kind of output
from `readpins`. It should look like this:

```
$ ~/go/src/github.com/abworrall/glassware/gw -v=1
 SerialController(C0), read(/dev/ttyUSB0)= "Controller:C0 A0:773 A1:673 A2:586 A3:507 A4:426 A5:383 "
 SerialController(C0), read(/dev/ttyUSB0)= "Controller:C0 A0:773 A1:674 A2:589 A3:511 A4:431 A5:388 "
 SerialController(C0), read(/dev/ttyUSB0)= "Controller:C0 A0:773 A1:674 A2:588 A3:509 A4:429 A5:384 "
```

Any unconnected analog pins will return noisy values.

## Ardiuno: FSRs

For each FSR, you should wire it up to an analog pin on the arduino
(using the breadboard) with a 10Kohm pull-down resistor, as per [this
guide](https://learn.adafruit.com/force-sensitive-resistor-fsr/using-an-fsr).
Using a pull-down means the voltage signal will increase as the
applied force increases.

TBD a photo of a setup with two FSRs

Once you have an FSR setup (say, on pin A0), run the `gw -v=1` command
again, to see the values being reported by the FSR. Squeeze the FSR as
hard as you can to see the max value, and then play with putting
different weights on top of it.

### FSRs are finicky

They are not accurate. You can put a known weight on/off the FSR a
bunch of times, and each time get a different answer. All you can get
reliably is a drop or jump in the right direction.

They are finicky in another way - the object's base needs to be
entirely within the sensor. If it is partly outside, then the amount
of weight borne by the FSR can vary enormously, and perhaps even be
zero. The simplest way to manage this is to insert a quarter coin
between the object and the FSR, so that all weight goes through the
quarter.

## Pi: host setup

Install some kind of Linux on the Pi. I went with
[raspbian](https://www.raspberrypi.com/software/). If Pies remain
weirdly unavailable, you can use a random linux box.

Run `sudo raspi-config` and setup the following:
- [System options > Wireless] - configure the Wifi your system will be using
- [System options > password] - for the non-root user
- [System options > hostname] - `glasspi` ?
- [System options > audio] - use headphone jack, not HDMI
- [System options > boot] - console, autologin
- [Interface options > SSH] - enable, so you can ssh/scp to the Pi

You can plug the Arduino board into the Pi, instead of your dev
machine, at this point.

You should also connect the headphone jack to your speaker, and check basic audio is working:
`aplay /usr/share/sounds/alsa/Front_Center.wav`.

Now setup the software on the Pi.

## Pi: the gw tool

From your dev machine, `scp ./gw pi@glasspi:~`, assuming `gw` was
cross-compiled for ARM as described above, and that your non-root user
is `pi`.

On the pi, you can test it by just running `~/gw -v` - it should
connect to the Arduino controller, and start printing out sensor
readings.

When you're happy it all works, you will want the tool to start on
boot, so:

TBD, auto-run tool in loop mode

## Pi: raspotify

This package will turn your Pi into a smart speaker that Spotify will
stream to, via Spotify Connect.

1. Set up raspotify

There is a [raspotify setup
guide](https://github.com/dtcooper/raspotify/wiki/Basic-Setup-Guide).
The 'Easy Way' doesn't work, because PulseAudio. The Hard Way ends up
needing you to do this, for a Raspbian Pi:

```
sudo cat <<EOT >> /etc/asound.conf
defaults.ctl.card 0
defaults.pcm.card 0
defaults.pcm.dmix.rate 44100
defaults.pcm.dmix.format S16_LE
EOT
```

Test it: `speaker-test -c2 -l1` should generate some nice pink noise.

2. Use your phone to link your new smart speaker into your Spotify account

Get your phone, and run the spotify app. Go to [Menu>Devices>Devices
Menu], and wait a little bit until you see something like `raspotify
(glasspi)` show up in the list. Then select it. This makes the Pi a
smart speaker destination that your Spotify account can play to, via
their Spotify Connect API.

This should be a one-time operation, but **you may need to repeat it
from time to time**, because Spotify Connect seems pretty flakey, and to
get confused when you use your spotify account to play music on your
phone or whatever.

### Pi: spotify login

These steps will authenticate the `gw` tool, and let it use the
Spotify web API to control streaming to your shiny new Pi-based smart
speaker.

This is kind of a PITA, but is a one-time setup operation.

1. Create yourself an 'app' on spotify
- go to https://developer.spotify.com/dashboard/, log in
   - should turn your account into a 'developer account' at some point
- create a new app (call it 'glassware' or whatever)
- edit settings, add a redirect URI: `http://localhost:8081/oauth-callback`
- get the app's `Client ID` and `Client Secret`, cut-n-paste 'em somewhere

2. Log the `gw` tool into spotify
- **The easy way:**
   - do it all on your dev machine, and copy the token over
   - run `gw -spotify-init -spotify-id=DEADBEEF -spotify-secret=DEADBEEF` (but using your ID and Secret)
   - after logging in, clicking 'agree', and being redirected back, the tool should say something like
   - ```
2022/12/01 14:34:14 Stored the OAuth2 token: /home/abw/.gw/spotify-oauth-token.json
2022/12/01 14:34:14 Have a spotify client logged in as: Adam
```
   - now copy that token over to your Pi: `scp -r ~/.gw/ pi@glasspi:~/`
- **The harder way":**
   - do it directly on the pi
   - you'll need a monitor/keyboard/mouse attached to your pi
   - if it is booting into a terminal, start up desktop by running `startx`
   - start a browser, open a terminal window
   - run `gw -spotify-init -spotify-id=DEADBEEF -spotify-secret=DEADBEEF` (but using your ID and Secret)

After this is complete, the `gw` tool can reuse the oauth2 token
indefinitely.

### Pi: bluetooth speakers

Having the Pi send audio to a bluetooth speaker is doable, but the
internet says you'll want a dedicated bluetooth dongle for the Pi
([link](https://www.sparkfun.com/products/17598)), because the Pi's
builtin bluetooth is flakey when the Pi is trying to do both Bluetooth
and Wifi, e.g. streaming music.

You'll need to figure out how to pair the speaker with the Pi, how to
reconfigure alsa (`/etc/asound.conf`), and repeat the "Hard Way" steps
in configuring raspotify.

## Spotify

Now the fun bit - deciding what music will be played when the sensors
detect weight changes !

Each sensor should get a corresponding playlist. The names are
predetermined and based on the sensor names, which in turn are based
on which arduino pins the sensors are connected to. If your sensor is
wired to pin `A3`, then it will try to play the playlist `Glassware
C0/A3`.
