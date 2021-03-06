package main

type Template struct {
	SuccessMsg string
	InputType  string
	Domain     string
	Path       string
	Port       string
}

func (self *Template) fallbackToFormInputTypeText() {
	// Some versions of android's default browser do
	// not handle <input type="time"> properly.
	self.InputType = "text"
}

var createReminderTemplate = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Send me a reminder</title>
    <meta name="description" content="Add reminders through a simple API ">
    <meta name="viewport" content="width=device-width, initial-scale=1">
  </head>
  <body>
    <h2>Send me a reminder!</h2>
    <p><strong>{{.SuccessMsg}}</strong></p>
    <form id="rem-form" action="{{.Domain}}:{{.Port}}/{{.Path}}" method="POST">
      <div>
        <input id="client-now" type="hidden" name="client-now">
      </div>
      <div>
        <label for="time">Time*: </label>
	<input id="time" type="{{.InputType}}" name="time" required>
      </div>
      <div>
        <label for="message">Message*: </label>
        <input type="text" name="message" required>
      </div>
      <div>
        <label for="date">Date: </label>
        <input id="date" type="date" name="date">
      </div>
      <div><input type="submit" value="Submit"></div>
    </form>
    <script>
      now = new Date();
      document.getElementById("client-now").setAttribute("value", now);
      then = new Date(now.getTime() + 30*60000);
      inThirtyMinutes = addZero(then.getHours()) + ":" + addZero(then.getMinutes());
      dateToday = now.getFullYear() + "-" + addZero(now.getMonth() + 1) + "-" + addZero(now.getDate());
      document.getElementById("time").setAttribute("value", inThirtyMinutes);
      document.getElementById("date").setAttribute("value", dateToday);
      function addZero(i) {
        if (i < 10) {
          i = "0" + i;
        }
        return i;
      }
      document.getElementById("rem-form").reset();
    </script>
  </body>
</html>
`
