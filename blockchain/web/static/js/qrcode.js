var qrcode = new QRCode(document.getElementById("qrcode"), {
    width : 300,
    height :300
});

function makeCode () {
    var elText = document.getElementById("text");

    if (!elText.value) {
        alert("Input a text");
        elText.focus();
        return;
    }

    qrcode.makeCode("http://172.19.172.147:9000/findTeaByID?id=" + elText.value);
}

makeCode();

$("#text").
on("blur", function () {
    makeCode();
}).
on("keydown", function (e) {
    if (e.keyCode == 13) {
        makeCode();
    }
});