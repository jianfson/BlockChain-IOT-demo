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

    qrcode.makeCode("http://223.128.94.60:9000/findTeaByID?id=" + elText.value);
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
