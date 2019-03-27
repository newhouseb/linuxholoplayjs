# Usage

From source
```
sudo go run main.go
```

From a release
```
sudo ./linuxholoplayjs
```

Navigate to a holoplay.js example and things should work!

# Notes

- This might work on other OSes (when you compile for them)!
- I've only tested this on one Looking Glass on one machine. Your mileage may vary.
- It currently presents the configuration to the holoplay.js SDK. I'm not sure if there are other features hidden in the SDK I don't know about.
- If you don't want to run it as root I think you can modify your udev rules as mentioned here: https://github.com/lonetech/LookingGlass/blob/master/20-lookingglass.rules

# Thanks

- I based a bunch of the USB HID protocol work based on sluething done by https://github.com/lonetech for https://github.com/lonetech/LookingGlass. I couldn't get it to work on my Linux install (because Python).

# License

- Because this links against a go hid package that in turn links against libusb, this is GNU LGPL 2.1 on linux. On other platforms that don't use libusb (i.e. Windows, Mac) it's 3-clause BSD. 
