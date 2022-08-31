# About

Some slightly customized Vagrantfiles for various use cases.


## Basic Vagrant usage / cheatsheet

Create default Vagrantfile:
```
vagrant init .
```

Booting up the Vagrant machine / creating it for the first time:
```
vagrant up
```

SSH into Vagrant machine (go to the directory where the `Vagrantfile` is located):
```
vagrant ssh
```

Get status of the Vagrant machine:
```
vagrant status
```

Powering down:
```
vagrant halt
```

Updating boxes
```
vagrant box update
```
