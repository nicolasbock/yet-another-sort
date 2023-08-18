#
# Regular cron jobs for the yet-another-sort package.
#
0 4	* * *	root	[ -x /usr/bin/yet-another-sort_maintenance ] && /usr/bin/yet-another-sort_maintenance
