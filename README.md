# Opal Activity

Simple script to export your Opal Card monthly activity to a CSV file https://www.opal.com.au/

I write this to find how many days I worked from home last financial year.

# Usage

```
go run main.go \ -username=your-username 
                 -password=your-password 
                 -month=month-of-activity 
                 -year=year-of-activity
```

### Example

```
go run main.go -username=rudylee -password=123456 -month=8 -year=2017
```
