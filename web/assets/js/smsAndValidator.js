var inter;
/**
 * 发送验证码
 */
function sendSMS() {
    //校验手机号，触发 bootstrapValidator 对手机号校验
    //初始化 bootstrapValidator 对象
    var validator = $("#myForm").data('bootstrapValidator');
    validator.validateField("user_tel");
    var flag = validator.isValidField("user_tel");
    console.log(flag);
    // 手机号被输入，且手机号规则符合正则才可以点击发送验证码
    if (flag) {
        //触发重复行为：每隔一秒显示一次数字
        inter = setInterval("showCount()", 1000);
        $(".qrcode").attr("disabled", true);

        // 请求服务器发送验证码
        $.post("${pageContext.request.contextPath}/sms", {
            "methodName": "sendSMS",
            "phoneNum": $("#user_tel").val()
        }, function (data) {
            console.log(data);
        }, "json");
    }
}

var count = 6; // 一般10或30或60s
function showCount() {
    $(".qrcode").text(count + "S");
    count--;
    if (count < 0) {
        clearInterval(inter);
        $(".qrcode").text("发送验证码");
        count = 6;
        $(".qrcode").attr("disabled", false);
    }
}

$(function () {
    $("#myForm").bootstrapValidator({
        message: "this is no a valiad field",
        fields: {
            user_tel: { // 手机号校验
                message: "手机号格式错误",
                validators: {
                    notEmpty: {
                        message: "手机号不能为空"
                    },
                    stringLength: {
                        message: "手机号长度为11",
                        min: 11,
                        max: 11
                    },
                    regexp: {
                        message: "手机号格式不对",
                        regexp: /^[1]{1}[1356789]{1}[0-9]+$/
                    }
                }
            },
            user_password: {
                message: "密码格式错误",
                validators: {
                    notEmpty: {
                        message: "密码不能为空"
                    },
                    stringLength: {
                        message: "密码长度为6~8",
                        min: 6,
                        max: 8
                    },
                    regexp: {
                        message: "密码由小写字母、数字组成",
                        regexp: /^[a-z0-9]+$/
                    },
                    different: {
                        message: "密码不能和手机号一致",
                        field: "user_tel"
                    }
                }
            },
            qrCode: { // 验证码输入框
                message: "验证码格式错误",
                validators: {
                    notEmpty: {
                        message: "验证码不能为空"
                    },
                    stringLength: {
                        message: "验证码长度为4",
                        min: 4,
                        max: 4
                    }
                }
            }
        }
    });
})
