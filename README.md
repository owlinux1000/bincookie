# bincookie

## What's this

A simple tool to parse `Cookies.binarycookies` which is often found in iOS.

## Motivation

A parsing tool for binarycookies has already published such as [BinaryCookieReader](https://github.com/as0ler/BinaryCookieReader). However, this tools is implemented by Python2. So, I want to replace it.

## Feature

- The output is able to use with `curl`

## How to use

Install to use `go get` or download from Releases.
```
$ go get github.com/owlinux1000/bincookie
$ bincookie Cookies.binarycookies
# Netscape HTTP Cookie File
# This file was generated by owlinux1000's bincookie
# https://github.com/owlinux1000/bincookie

.example.com    TRUE    /       FALSE   1550464816      session     hogehogehogehoge
```
