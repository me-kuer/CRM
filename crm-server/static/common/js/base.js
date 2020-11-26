// 主机
var host = "";

// 登录页面的路由
var login_router = "/login";



/**
 *
 * @param params
 * @returns {Promise<any>}
 */
function ajax(params){
    //接收参数
    var request = {};
    request.url = params.url;
    request.type = params.type ? params.type : "get";
    request.is_formdata = params.is_formdata ? params.is_formdata : false;
    request.data = params.data ? params.data : null;
    var promise = new Promise(function(resolve, reject){
        var xhr = new XMLHttpRequest();
        xhr.onreadystatechange = function(){
            if(xhr.readyState == 4){
                var res = JSON.parse( xhr.responseText );
                // 如果返回的错误码为401，则直接跳转到登录页面
                if(res.code == 401){
                    top.location.href = login_router;
                }
                resolve(res);
            }
        }
        xhr.open(request.type, request.url);

        // 判断提交的数据是否为formData类型
        if(request.is_formdata){
            xhr.send(request.data);
        }else{
            xhr.setRequestHeader ("content-type", "application/x-www-form-urlencoded" );
            xhr.send(json2url(request.data));
        }
    });
    return promise;
}

/**
 * @param str
 * @returns {boolean}
 */
function isNull( str ){
    if ( str == "" ) return true;
    var regu = "^[ ]+$";
    var re = new RegExp(regu);
    return re.test(str);
}

/**
 * json对象 转 url参数
 * @param {json} data
 */
function json2url(data) {
    var tempArr = [];
    for (var i in data) {
        var key = encodeURIComponent(i);
        var value = encodeURIComponent(data[i]);
        tempArr.push(key + "=" + value);
    }
    var urlParamsStr = tempArr.join("&");
    return urlParamsStr;
}

//删除字符左右两端的空格
function trim(str){
    return str.replace(/(^\s*)|(\s*$)/g, "");
}

/**
 * 获取url中的参数值
 * @param {string} name 参数名称
 */
function getUrlParam(name)
{
    var query = window.location.search.substring(1);
    var vars = query.split("&");
    for (var i=0;i<vars.length;i++) {
        var pair = vars[i].split("=");
        if(pair[0] == name){return pair[1];}
    }
    return(false);
}
