var websocket = null;
var wsID = null;
var authenticated = false;


function promptPassword() {
    return prompt("Please enter tonight's meeting password.");
}

function auth() {
    if (!authenticated) {
        let password = promptPassword();

        // synchronous request to /authenticate
        var xhttp = new XMLHttpRequest();
        xhttp.open("POST", "/authenticate", false);
        xhttp.setRequestHeader("Content-Type", "application/json");

        let event = {
            "id": wsID,
            "password": password
        };

        xhttp.send(JSON.stringify(event));

        authenticated = xhttp.status == 200;
    }

    return authenticated;
}