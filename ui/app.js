//############################# Port Check ###########################################//
function openPortCheck() {
    $('#portCheckModal').modal('show');
}

function submitForm() {
    var formData = new FormData(document.getElementById('portCheckForm'));
    var jsonData = {};

    formData.forEach(function (value, key) {
        jsonData[key] = value;
    });

    fetch('/port', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(jsonData),
    })
        .then(response => response.json())
        .then(data => {
            var outcomeContainer = document.getElementById('outcomeContainer');
            if (data.status === "Connected") {
                outcomeContainer.innerHTML = '<div class="alert alert-success" role="alert">' + data.status + '</div>';
                outcomeContainer.style.backgroundColor = 'green'; // Set background color to green
            } else {
                outcomeContainer.innerHTML = '<div class="alert alert-danger" role="alert">' + data.errorMessage + '</div>';
                outcomeContainer.style.backgroundColor = 'red'; // Set background color to red
            }
        })
        .catch(error => console.error('Error:', error));

    // Do not close the modal here
}

function closeModal() {
    $('#portCheckModal').modal('hide');
}

//############################# Domain Check ###########################################//

document.getElementById('toggleDnsInput').addEventListener('change', function () {
    var dnsInput = document.getElementById('dnsInput');
    dnsInput.style.display = this.checked ? 'block' : 'none';
});

function resolveDomain() {
    $('#domainCheckModal').modal('show');
}

function closeDomainModal() {
    $('#domainCheckModal').modal('hide');
}

function submitDomainCheck() {
    var jsonData = {};
    var domainName = document.getElementById('domainName').value;
    var dnsServer = document.getElementById('dnsServer').value;

    if (dnsInput.style.display === 'none') {
        dnsServer = ""
    }

    jsonData['domainName'] = domainName;
    jsonData['dnsServer'] = dnsServer;

    fetch('/resolve', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(jsonData),
    })
        .then(response => response.json())
        .then(data => {
            // Display domain resolution result if no error message
            var domainResult = document.getElementById('domainResult');
            if (data.status === "Error") {
                domainResult.innerHTML = '<div class="alert alert-danger" role="alert">' + data.errorMessage + '</div>';
            } else {
                domainResult.innerHTML = `
                            <h5>Domain Name: ${data.domainName}</h5>
                            <table class="table">
                            <thead>
                                <tr>
                                <th scope="col">Type</th>
                                <th scope="col">Records</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr>
                                <td>A</td>
                                <td>${data.A.join(', ')}</td>
                                </tr>
                                <tr>
                                <td>CNAME</td>
                                <td>${data.CNAME ? data.CNAME.join(', ') : ''}</td>
                                </tr>
                                <tr>
                                <td>MX</td>
                                <td>${data.MX.join(', ')}</td>
                                </tr>
                            </tbody>
                            </table>
                            <p>Response Time: ${data.responseTime} ms</p>
                            <p>DNS Server: ${data.dnsServer}</p>
                        `;
            }
        })
        .catch(error => {
            console.error('Error:', error);
            // Handle errors here
        });
}

//############################# Http Check ###########################################//
function openHttpModal() {
    $('#httpCheckModal').modal('show');
}

function closeHttpModal() {
    $('#httpCheckModal').modal('hide');
}

function submitHttpCheck() {
    var url = document.getElementById('url').value;
    var jsonData = {}
    jsonData['url'] = url;

    fetch('/http', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(jsonData),
    })
        .then(response => response.json())
        .then(data => {
            var httpResult = document.getElementById('httpResult');
            httpResult.innerHTML = `
              <table class="table">
                <thead>
                  <tr>
                    <th>URL</th>
                    <th>Status</th>
                    <th>Status Code</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>${data.url}</td>
                    <td>${data.success ? 'Success' : 'Failure'}</td>
                    <td>${data.statusCode}</td>
                  </tr>
                </tbody>
              </table>
              <p>Response Time: ${data.responseTime} ms</p>
            `;
        })
        .catch(error => {
            console.error('Error:', error);
            // Handle errors here
        });
}


//############################# DB Check ###########################################//

function showCheckDBModal() {
    $('#dbCheckModal').modal('show');
}

function closeDBModal() {
    $('#dbCheckModal').modal('hide');
}

function submitDBForm() {
    var outcomeContainer = document.getElementById('dbCheckResult');
    outcomeContainer.innerHTML = '<div></div>';
    var dbType = document.getElementById("dbType").value;
    var formData = new FormData(document.getElementById('dbCheckForm'));

    var jsonData = {};
    formData.forEach(function (value, key) {
        jsonData[key] = value;
    });

    if (dbType === "oracle") {
        jsonData['serviceName'] = document.getElementById("serviceName").value;
    } else if (dbType === "postgres") {
        jsonData['databaseName'] = document.getElementById("databaseName").value;
    }

    fetch('/db', {
        method: 'POST',
        headers: {
            userdata: jsonToBase64(jsonData)
        }
    })
        .then(response => response.json())
        .then(data => {
            var outcomeContainer = document.getElementById('dbCheckResult');
            if (data.status === "Success") {
                outcomeContainer.innerHTML = '<div class="alert alert-success" role="alert">' + data.status + '</div>';
                outcomeContainer.style.backgroundColor = 'green'; // Set background color to green
            } else {
                outcomeContainer.innerHTML = '<div class="alert alert-danger" role="alert">' + data.errorMessage + '</div>';
                outcomeContainer.style.backgroundColor = 'red'; // Set background color to red
            }
        })
        .catch(error => console.error('Error:', error + response));

    // Do not close the modal here
}

function jsonToBase64(object) {
    const json = JSON.stringify(object);
    var utf8String = unescape(encodeURIComponent(json));
    return btoa(utf8String);
}

function changeDbType() {
    var dbType = document.getElementById("dbType").value;

    if (dbType === "oracle") {
        document.getElementById("serviceNameDiv").style.display = "block";
        document.getElementById("dataBaseDiv").style.display = "none";
    } else if (dbType === "postgres") {
        document.getElementById("serviceNameDiv").style.display = "none";
        document.getElementById("dataBaseDiv").style.display = "block";
    } else {
        document.getElementById("serviceNameDiv").style.display = "none";
        document.getElementById("dataBaseDiv").style.display = "none";
    }
} 