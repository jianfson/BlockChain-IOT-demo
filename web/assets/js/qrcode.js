var qrcode = new QRCode(document.getElementById("qrcode"), {
    width : 540,
    height :540
});

function makeCode () {
    var elText = document.getElementById("text");

    if (!elText.value) {
        alert("Input a text");
        elText.focus();
        return;
    }
    qrcode.makeCode("http://47.108.134.136:9000/findTeaByID?id=" + elText.value);
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
