var oInput=document.getElementById('popout_button');
var box=document.getElementById('popout');
var a=true;
oInput.onclick=function(){
    a=!a;
    a?box.style.display="block":box.style.display="none";
}

