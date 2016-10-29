# Journalcheck – like logcheck, but for journald

Journalcheck is a simple utility which sends mails of possibly
interesting [journal] entries to a configurable e-mail address.
Journal entries are marked as not interesting by matching a regular
expression, a list of which can be provided in a file, one per line.

This is very similar in operation to [logcheck], except it operates on
the journal instead of plain text log files.

[journal]: https://www.freedesktop.org/software/systemd/man/systemd-journald.service.html
[logcheck]: http://logcheck.org/

## Installation

```
go get github.com/jorgenschaefer/journalcheck
mkdir -p /opt/journalcheck/{bin,etc,var}
cp $GOPATH/bin/journalcheck /opt/journalcheck/bin/
touch /opt/journalcheck/etc/filters.txt
cp $GOPATH/src/github.com/jorgenschaefer/journalcheck/default /etc/default/journalcheck
cp $GOPATH/src/github.com/jorgenschaefer/journalcheck/journalcheck.{service,timer} /etc/systemd/system
```

Now edit `/etc/default/journalcheck` and
`/opt/journalcheck/etc/filters.txt`. The latter should contain regular
expressions matching the lines `journalcheck -o match` emits. You can
use that command in conjunction with `-n 1000` or similar to get a
larger number of lines you might want to ignore.

Journalcheck sends mail using `/usr/sbin/sendmail`, so make sure that
works (using e.g. [nullmailer](http://untroubled.org/nullmailer/)).

## Configuration

Journalcheck is configured using either command line arguments (see
`journalcheck -h`) and/or the following environment variables:

- `JOURNALCHECK_CURSORFILE`: A file to store the last read cursor
  in. Needs to be writable by the journalcheck process.
- `JOURNALCHECK_FILTERFILE`: A file containing regular expressions,
  one per line, matching entries to ignore. These regular expressions
  are matched against `<identifier>: <message>` lines, which you can
  see using the `-o match` format argument. This is different from the
  default format, as it does not include time stamps, the host name,
  or the PID of the process.
- `JOURNALCHECK_RECIPIENT`: An e-mail address to send mails to.

