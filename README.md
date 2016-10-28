# journalcheck – like logcheck, but for journald

journalcheck is a simple utility which takes a whitelist file of
regular expressions, one per line, and then follows the journal.
Entries that do not match any of the regular expressions will be sent
to an e-mail address as potentially interesting.

This is very similar to [logcheck](http://logcheck.org/).
