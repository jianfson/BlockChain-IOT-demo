var aaa = document.getElementById("btnAdd"); //��ȡ��ʾͼƬ��divԪ��
var input = document.getElementById("file_input"); //��ȡѡ��ͼƬ��inputԪ��
var tdid = 1;
//������жϱ�������Ƿ�֧�����API��
if(typeof FileReader==='undefined'){
    aaa.innerHTML = "��Ǹ������������֧�� FileReader";
    input.setAttribute('disabled','disabled');
}else{
    input.addEventListener('change',readFile,false); //���֧�־ͼ����ı��¼���һ���ı��˾�����readFile������
}


function readFile(){
    for (var index = 0; index<this.files.length; index++){
        var file = this.files[index]; //��ȡfile����
        //�ж�file�������ǲ���ͼƬ���͡�
        if(!/image\/\w+/.test(file.type)){
            alert("�ļ�����ΪͼƬ��");
            return false;
        }
    }
    for (var index = 0; index<this.files.length; index++){
        var file = this.files[index]; //��ȡfile����

        var reader = new FileReader(); //����һ��FileReaderʵ��
        reader.readAsDataURL(file); //����readAsDataURL��������ȡѡ�е�ͼ���ļ�
        //�����onload�¼��У���ȡ���ɹ���ȡ���ļ����ݣ����Բ���һ��img�ڵ�ķ�ʽ��ʾѡ�е�ͼƬ
        reader.onload = function(e){
            tdid++;
            $('<li id=' + tdid + ' style="position:relative;"><div class="imgwrap">'
                + '<img src="../imgs/topic_pic_def.png" alt=""/><div class="bar mint active" style="position:absolute;width:80%;height:0.6rem;top:5px;left:5px;right:5px;" data-percent="100" ></div></div></li>').insertBefore($("#tdAdd"));
            var imageStr = this.result;
            /*�ӳ���ʾͼƬ ģ��ͼƬ�ϴ��ɹ������ʾ��
             �������ֱ����ʾͼƬ�ˡ�����Ȥ��ͬѧ����ʵ���ϴ��еĽ�����Ч������
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

            /* �ϴ�ͼƬ����̨���ز���ʾ��
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