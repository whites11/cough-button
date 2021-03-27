# Cough button

This project aims at providing software support for a cheap DYI cough button for linux.

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

- looks for a supported serial device to listen inputs from and to port status to
- connects to pulseaudio via DBUS, setting the mute status and getting notifications about state changes.

The LED of the external device is updated to reflect in near-realtime the muted status of the default source of the system.

WARNING: if you are using custom settings for some applications, `cough button` will not be able to mute those applications for you. It always mutes
the default source.
