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
- install the Go programming language, [golang](https://go.dev/), setting up `$GOPATH`, `~/go`, etc.

## Get and build the glassware software

```
# Check out the code, wheren golang toolchain expects to see it
mkdir -p ~/go/src/github.com/abworrall
cd  ~/go/src/github.com/abworrall
git clone https://github.com/abworrall/glassware
cd glassware

# Compile the command, so you can run it locally
go build ./cmd/gw
```

To cross-compile the command for a Raspberry Pi, set a few ENV vars:
```
GOOS=linux GOARCH=arm go build ./cmd/gw
```

## The Raspberry Pi host

Install some kind of Linux on the Pi. I went with
[raspbian](https://www.raspberrypi.com/software/). If Pis remain
weirdly unavailable, you can use a random linux box.

You'll need to enable SSH services on the Pi, it is disabled by default.

## The arduino controller board

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
... SerialController(C0), read(/dev/ttyUSB0)= "Controller:C0 A0:773 A1:673 A2:586 A3:507 A4:426 A5:383 "
... SerialController(C0), read(/dev/ttyUSB0)= "Controller:C0 A0:773 A1:674 A2:589 A3:511 A4:431 A5:388 "
... SerialController(C0), read(/dev/ttyUSB0)= "Controller:C0 A0:773 A1:674 A2:588 A3:509 A4:429 A5:384 "
```

Any unconnected analog pins will return noisy values.

## The FSRs

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

## Software setup on the Pi

1. Put `gw` on the Pi, set it up to run on boot
2. double check perms on the pi, for user access to /dev/ttyUSB0 etc
3. install mopidy, configure it for Spotify
4. configure mopidy to talk to your speaker (or use aux out)
5. on Spotify setup a few playlists (named for sensors: "Glassware C0/A0", etc)
