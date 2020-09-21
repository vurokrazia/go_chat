document.addEventListener('DOMContentLoaded', function () {
  console.log("grettings")
});

function mySubmitFunction(e) {
  e.preventDefault();
  send_validation()
  return false;
}

function send_validation() {
  fetch('http://localhost:8000/validate',
  {
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json'
    },
    method: 'POST',
    body: JSON.stringify({
      user_name: user_name
    }),
  })
  .then(response => response.json())
  .then(data => console.log(data));
}

function success_result(data) {
  if (data.is_valid === true){
    create_conexion(data)
  } else {
    alert("Loooseeeer")
  }
}

function create_conexion(data) {
  var conexion = new WebSocket("ws://localhost:8000/chat" + data.user_name)
}