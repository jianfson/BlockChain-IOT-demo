var aaa = document.getElementById("btnAdd"); //获取显示图片的div元素
var input = document.getElementById("file_input"); //获取选择图片的input元素
var tdid = 1;
//这边是判断本浏览器是否支持这个API。
if(typeof FileReader==='undefined'){
    aaa.innerHTML = "抱歉，你的浏览器不支持 FileReader";
    input.setAttribute('disabled','disabled');
}else{
    input.addEventListener('change',readFile,false); //如果支持就监听改变事件，一旦改变了就运行readFile函数。
}


function readFile(){
    for (var index = 0; index<this.files.length; index++){
        var file = this.files[index]; //获取file对象
        //判断file的类型是不是图片类型。
        if(!/image\/\w+/.test(file.type)){
            alert("文件必须为图片！");
            return false;
        }
    }
    for (var index = 0; index<this.files.length; index++){
        var file = this.files[index]; //获取file对象

        var reader = new FileReader(); //声明一个FileReader实例
        reader.readAsDataURL(file); //调用readAsDataURL方法来读取选中的图像文件
        //最后在onload事件中，获取到成功读取的文件内容，并以插入一个img节点的方式显示选中的图片
        reader.onload = function(e){
            tdid++;
            $('<li id=' + tdid + ' style="position:relative;"><div class="imgwrap">'
                + '<img src="../imgs/topic_pic_def.png" alt=""/><div class="bar mint active" style="position:absolute;width:80%;height:0.6rem;top:5px;left:5px;right:5px;" data-percent="100" ></div></div></li>').insertBefore($("#tdAdd"));
            var imageStr = this.result;
            /*延迟显示图片 模拟图片上传成功后的显示。
             我这里就直接显示图片了。有兴趣的同学可以实现上传中的进度条效果。。
           */
            setTimeout(function(){
                var td = $("#" + tdid);
                td.html("<div class='imgwrap'><img src='" + e.target.result + "'/></div>");
                var $closeImg = $('<img src="../imgs/close_btn.png" style="position:absolute;top:5px;right:5px;width:20px;height:20px;">').appendTo(td);
                //  $('<input type="hidden" name="imguuids" class="imgHidden" value="'+data.uuid+'"/>').appendTo(td);
                $closeImg.click(function(){
                    $(this).closest("li").remove();
                });
            }, 2000);

            /* 上传图片到后台返回并显示。
            $.ajax({
                url: "upload.jhtm",
                type: "POST",
                data: {tdid :tdid, imageStr:imageStr},
                dataType: "json",
                cache: false,
                success: function(data) {
                    var rtntdid = data.tdid;
                    var td = $("#" + rtntdid);
                    td.html("<div class='imgwrap'><img src='" + data.medium + "'/></div>");
                    var $closeImg = $('<img src="../imgs/close_btn.png" style="position:absolute;top:5px;right:5px;width:20px;height:20px;">').appendTo(td);
                    $('<input type="hidden" name="imguuids" class="imgHidden" value="'+data.uuid+'"/>').appendTo(td);
                    $closeImg.click(function(){
                        $(this).closest("td").remove();
                    });
                }
            });
            */
        }
    }
}