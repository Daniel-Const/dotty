# Dotty

A simple Dot file manager!

- Manage dot files in one git repository
- Create multiple profiles of dotfile configs
- Easily deploy dotfile config on any machine

## Usage

### Install

```bash
> git clone https://github.com/Daniel-Const/dotty.git
> cd dotty
> mkdir bin/ && go build -o bin/
```

## Problem

I have multiple devices with different operating systems. I want to be able to configure them with different dotfile setups,
manage them in one repository, and easily deploy / update them on different machines.

## Solution

A Go CLI app that automatically copies all of the dot files to the right locations :)

- The dotty.map file maps dotfiles in a profile to deploy locations.
- This way I can create profiles for my pc, laptop etc. Manage them in one repository, and easily deploy the configs to each machine

## Plan

![alt text](./public/dotty-plan.png)
