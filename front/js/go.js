document.addEventListener('DOMContentLoaded', function () {
  console.log("grettings")
});

var user_name;
var last_conexion;
var last_message = "";
var button_disabled = false

function mySubmitFunction(e) {
  e.preventDefault();
  send_validation(e.target[0].value)
  return false;
}

function send_validation(name) {

  fetch('http://localhost:8000/validate',
    {
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      },
      method: 'POST',
      body: JSON.stringify({
        user_name: name
      }),
    })
    .then(response => response.json())
    .then(data => {
      if (data.IsValid) {
        document.getElementById("nickname").style.display = "none"
        document.getElementById("container_chat").style.display = "block"
        user_name = name
        success_result(data)
      }
    });
}

function success_result(data) {
  if (data.IsValid === true) {
    create_conexion()
  } else {
    alert("Loooseeeer")
  }
}

function send_msg() {
  if (button_disabled) return
  last_conexion.send(document.getElementById("msg").value)
  document.getElementById("msg").value = ""
}

function create_conexion() {
  var conexion = new WebSocket("ws://localhost:8000/chat/" + user_name)
  last_conexion = conexion
  conexion.onopen = function (t) {
    conexion.onmessage = function (respose) {
      if (last_message != respose.data) {
        var values = document.getElementById("chat_area").value + "\n"
        document.getElementById("chat_area").value = values + respose.data
        last_message = respose.data
        setTimeout(() => {
          last_message = ""
          button_disabled = false
        }, 1000);
      }
    }
  }
}