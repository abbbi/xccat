# xccat
Shows latest [DHV-XC](https://www.dhv-xc.de)  flights on console

 * "top" like refresh interval
 *  Colored or ascii only output
 *  Limit results
 *  Filter results based on XC points
 *  Filter results based on takeoff (for example -f Wallberg)

See also:

 * [xcup](https://github.com/abbbi/xcup) for publishing flights.
 * [xcbackup](https://github.com/abbbi/xcbackup) for exporting flights.

![Alt text](xccat.jpg?raw=true "Screenshot")

# usage
```
Usage:
  xccat [OPTIONS]

Application Options:
  -d, --day=      date selection: 08.06.2022
  -i, --interval= Refresh interval in seconds (default: 0)
  -l, --limit=    Limit to X results (default: 0)
  -p, --points=   Only show flights >= XC points (default: 0)
  -a, --ascii     Dont display colors, ascii only output
  -f, --takeoff=  Filter by takeoff: takeoff must include string
  
Help Options:
  -h, --help      Show this help message
```
