# Cough button

This project aims at providing software support for a cheap (normally < 10€) DYI hardware cough button for linux.

A cough button is simply a button that allows to mute and unmute a microphone attached to a linux computer.

Theoretically the computer-side software provided with this repo can be easily extended/replaced to support other operating systems as well. 
I don't care about any other operating system so I will probably never do it (but feel free to contribute).

## Diagram of the components

![Diagram of the components](https://github.com/whites11/cough-button/blob/master/schema.png?raw=true)

## Hardware

The idea of this project is to act as an easy-to-hack framework and be adaptable to different hardware and potentially different sound systems.
Nevertheless it is being developed and (poorly) tested with the following hardware:

- The sketch included in this repo is meant to be flashed on a teensy 2.0 device: https://www.pjrc.com/store/teensy.html
- The button I am personally using is [this one](https://it.aliexpress.com/item/1005001995455752.html?spm=a2g0o.productlist.0.0.220c9418EKjHbK&algo_pvid=745e5556-c755-401f-9dc3-5d9cbd19a814&algo_expid=745e5556-c755-401f-9dc3-5d9cbd19a814-14&btsid=2100bddd16168746826697861e0f11&ws_ab_test=searchweb0_0,searchweb201602_,searchweb201603_) but virtually any button can do.

This is the connections diagram:

![Connections diagram](https://github.com/whites11/cough-button/blob/master/diagram.png?raw=true)

## Concepts

The daemon is a go CLI tool that has to be run from the command line as unprivileged user.

During startup, the deamon:

- looks for a supported serial device to listen inputs from and to post the microphone muted status to.
- connects to pulseaudio via DBUS, setting the mute status and getting notifications about state changes.

The LED of the external device is updated to reflect in near-realtime the muted status of the default source of the system.

WARNING: if you are using custom settings for some applications, `cough button` will not be able to mute those applications for you. It always mutes
the default source.

## Getting started

It is still early days for the project, so running it is a manual process.

Hopefully my interest will in the project will continue and/or contributors will join me making it more straightforward.

In any case these are the steps needed to get yourself a cough button.

- Buy the hardware and connect following the connections diagram above.
- Flash the teensy device following [the official instructions](https://www.pjrc.com/teensy/teensyduino.html).
- Run the daemon as unprivileged user (cd into the `daemon` directory and run `go run main.go`).
