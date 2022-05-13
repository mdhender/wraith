# wraith
Wraith Game Engine

# Installation Notes
Please read the
[Systems Operator Manual](https://github.com/mdhender/wraith/blob/main/docs/sysop.adoc)
for instructions on installing and configuring the Wraith server.

# Player Notes
Please read the
[rulebook](https://github.com/mdhender/wraith/blob/main/docs/rulebook.adoc)
for instructions on playing Wraith.

# Influences
Wraith draws inspiration from Far Horizons, Empyrean Challenge, and The Campaign for North Africa.

# Sources
Logic for binding viper and cobra taken from
[Sting of The Viper](https://carolynvanslyck.com/blog/2020/08/sting-of-the-viper/).

Mat Ryer's Way sourced from
[matryer/way](https://github.com/matryer/way/commit/9632d0c407b008073d19d0c4da1e0fc3e9477508).

Start the server with the ability to shut it down gracefully from
[clavinejune blog](https://clavinjune.dev/en/blogs/golang-http-server-graceful-shutdown/).

See
[Gregory Gaines' Blog](https://www.gregorygaines.com/blog/how-to-properly-hash-and-salt-passwords-in-golang-bcrypt/)
for details on why we want to use BCrypt.

See
[Scott Piper's Blog](http://0xdabbad00.com/2015/04/23/password_authentication_for_go_web_servers/)
for more details on auth/auth.

## systemd
See the
[DO Tutorial](https://www.digitalocean.com/community/tutorials/how-to-sandbox-processes-with-systemd-on-ubuntu-20-04)
for details on securing and locking down this as a service.

FWIW, this is my starter:

    /etc/systemd/system# cat wraith.service
    [Unit]
    Description=Wraith API server
    StartLimitIntervalSec=0
    After=network-online.target
    
    [Service]
    Type=simple
    User=www-data
    PIDFile=/run/wraith.pid
    WorkingDirectory=/var/www/wraith
    ExecStart=/usr/local/bin/wraith
    ExecReload=/bin/kill -USR1 $MAINPID
    Restart=on-failure
    RestartSec=1
    
    [Install]
    WantedBy=multi-user.target
