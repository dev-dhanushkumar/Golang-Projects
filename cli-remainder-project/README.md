<center><h2> Golang CLI Task Remainder Tool</h2> </center>

<p style="text-align: justify">I'd be happy to help you write a description for your Golang CLI task reminder project. To make sure I provide the most relevant description, could you please share a few more details about the project? </p>
<b>For example: </b>
<li>What are the main features of your task reminder?</li>
<li>How does the user interact with it (command-line commands, flags, etc.)?</li>
<li>Does it have any special functionality like notifications, task prioritization, or recurring tasks?</li>
<li>Are there any unique aspects or goals for this project?</li>
<br>
<p style="text-align: justify">Once you provide some additional information, I can draft a concise project description for you. If you'd like a more detailed description afterwards, just let me know and I can expand on it.</p>

## [When](https://github.com/olebedev/when.git)
```when``` is a natural language date/time parser with pluggable rules and merge strategies.

Examples
* tonight at 11:10 pm
* at Friday afternoon
* the deadline is next tuesday 14:00
* drop me a line next wednesday at 2:25 p.m
* it could be done at 11 am past tuesday

<h4>How it works </h4>
<p style="text-align: justify">Usually, there are several rules added to the parser's instance for checking. Each rule has its own borders - length and offset in provided string. Meanwhile, each rule yields only the first match over the string. So, the library checks all the rules and extracts a cluster of matched rules which have distance between each other less or equal to options.Distance, which is 5 by default. For example:</p>

```code
on next wednesday at 2:25 p.m.
   └──────┬─────┘    └───┬───┘
       weekday      hour + minute
```
<h4>Usage</h4>

```go
w := when.New(nil)
w.Add(en.All...)
w.Add(common.All...)

text := "drop me a line in next wednesday at 2:25 p.m"
r, err := w.Parse(text, time.Now())
if err != nil {
	// an error has occurred
}
if  r == nil {
 	// no matches found
}

fmt.Println(
	"the time",
	r.Time.String(),
	"mentioned in",
	text[r.Index:r.Index+len(r.Text)],
)
```
## [Beeep](https://github.com/gen2brain/beeep)
```beeep``` provides a cross-platform library for sending desktop notifications, alerts and beeps.

<h4>Build Tags</h4>

```nodbus```- disable ```godbus/dbus``` and only ```notify-send```
<h4>Examples</h4>

```go
err := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
if err != nil {
    panic(err)
}
```
```go
err := beeep.Notify("Title", "Message body", "assets/information.png")
if err != nil {
    panic(err)
}
```
```go
err := beeep.Alert("Title", "Message body", "assets/warning.png")
if err != nil {
    panic(err)
}
```

## Install and run
Download and Install the github dependancy
```bash
$ go mod tiny
```

Run Project
```bash
$ go run main.go <hh:mm> <text message>
```
Example
```bash
$ go run main.go 23:45 "Time to go Sleep"
```