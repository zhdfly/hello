
//初始化计数器
var num = 0;
//区块锁定标识
var lock = false;
//加载layer拓展
layer.config({
    extend: 'extend/layer.ext.js'
});
//右键菜单参数
context.init({
    fadeSpeed: 100,
    filter: function ($obj) { },
    above: 'auto',
    preventDoubleContext: true,
    compress: false
});
function createBox(data) {
    var dataId = data.id || '';
    var value = data.text || '';
    var color = data.color || '';
    var height = data.height || 0;
    var width = data.width || 0;
    var pageX = data.pageX || 0;
    var pageY = data.pageY || 0;
    var tip = data.tip ||0;
    //更新计数器并记录当前计数
    var curNum = num++;
    //创建区域块
    var pos = $("#dragroom").position();
    var boxs = $('<div class="boxs"  rel="' + curNum + '" dataId="' + dataId + '"><div class="contentpic">' + value + '</div><div class="bg transparent" style="background-color:' + color + '"></div><div class="coors"></div><div class="tip">'+tip+'</div></div>').css({
        width: width,
        height: height,
        top: pageY > 0 ? pageY : (pos.top > 0 ? 0 : pos.top * -1 + 50),
        left: pageX > 0 ? pageX : (pos.left > 0 ? 0 : pos.left * -1 + 30)
    }).appendTo("#dragroom");
    console.log("log")
    var text_w = $(boxs.find('.contentpic')).width()
    var text_h = $(boxs.find('.contentpic')).height()
    text_w = boxs.find('.contentpic').parent().width() / 2 - text_w / 2
    text_h = boxs.find('.boxs').height() / 3 * -1 - text_h / 2
    console.log(text_w, text_h)
    //计算文本位置
    boxs.find('.contentpic').css({
        marginLeft: text_w,
        marginTop: text_h
    });
    //创建右键菜单
    context.attach('.boxs[rel=' + curNum + ']', [
        {
            text: '删除区域', action: function (e) {
                e.preventDefault();
                $('.boxs[rel=' + curNum + ']').remove();
            }
        },
        { divider: true },
        { header: '更改背景颜色' },
        {
            text: '<font color="#ef4836">红色</font>', action: function (e) {
                e.preventDefault();
                $('.boxs[rel=' + curNum + '] .bg').css('background-color', '#ef4836');
            }
        },
        {
            text: '<font color="#399bff">蓝色</font>', action: function (e) {
                e.preventDefault();
                $('.boxs[rel=' + curNum + '] .bg').css('background-color', '#399bff');
            }
        },
        {
            text: '<font color="#26a65b">绿色</font>', action: function (e) {
                e.preventDefault();
                $('.boxs[rel=' + curNum + '] .bg').css('background-color', '#26a65b');
            }
        },
        { divider: true },
        { header: '更改前景颜色' },
        {
            text: '<font color="#000">黑色</font>', action: function (e) {
                e.preventDefault();
                $('.boxs[rel=' + curNum + '] .contentpic').css('color', '#000');
            }
        },
        {
            text: '<font color="#000">白色</font>', action: function (e) {
                e.preventDefault();
                $('.boxs[rel=' + curNum + '] .contentpic').css('color', '#fff');
            }
        },
    ]);
}
function createvideoBox(data) {
    var dataId = data.id || '';
    var value = data.text || '';
    var color = data.color || '';
    var height = data.height || 0;
    var width = data.width || 0;
    var pageX = data.pageX || 0;
    var pageY = data.pageY || 0;

    //更新计数器并记录当前计数
    var curNum = num++;
    //创建区域块
    var pos = $("#dragroom").position();
    var boxs = $('<div class="top-label boxs" rel="' + curNum + '" dataId="' + dataId + '"><video id="myPlayer" width=100% height=100% class="video-js vjs-default-skin" controls autoplay><source src="' + value + '" type="application/x-mpegURL"></video><div class="coors transparent"></div></div>').css({
        width: width,
        height: height,
        top: pageY > 0 ? pageY : (pos.top > 0 ? 0 : pos.top * -1 + 50),
        left: pageX > 0 ? pageX : (pos.left > 0 ? 0 : pos.left * -1 + 30)
    }).appendTo("#dragroom");

    //计算文本位置
    boxs.find('.contentpic').css({
        marginLeft: boxs.find('.contentpic').width() / 2 * -1,
        marginTop: boxs.find('.contentpic').height() / 2 * -1
    });
    //创建右键菜单
    context.attach('.boxs[rel=' + curNum + ']', [
        {
            text: '删除区域', action: function (e) {
                e.preventDefault();
                $('.boxs[rel=' + curNum + ']').remove();
            }
        },
    ]);
}
//添加区域
$("#btn_add").click(function () {
    //弹出区域说明输入框

    createBox({
        text: $("#dotname").val(),
        width: 100,
        height: 100
    });

});
$("#btn_addvideo").click(function () {
    //弹出区域说明输入框
    createvideoBox({
        text: "http://hls.open.ys7.com/openlive/f01018a141094b7fa138b9d0b856507b.m3u8",
        width: 100,
        height: 100
    });
    new EZUIPlayer('myPlayer');
});
//添加区域
$("#btn_gd").click(function () {
    //弹出区域说明输入框
    layer.prompt({
        title: '请输入区域说明',
        formType: 0 //0:input,1:password,2:textarea
    }, function (value, index, elem) {
        layer.close(index);
        creategd({
            text: value,
            width: 100,
            height: 100
        });
    });
});
//锁定区域
$('#btn_lock').click(function () {
    if (lock) {
        $(this).val("锁定区域");
        lock = false;
        $('.boxs .coors').show();
    } else {
        $(this).val("解锁区域");
        lock = true;
        $('.boxs .coors').hide();
    }
});
//获取所有区块
$('#btn_save').click(function () {
    var data = [];
    $('.boxs').each(function () {
        var boxs = {};
        boxs['id'] = $(this).attr('dataId');
        boxs['tip'] = $(this).attr('tip');
        boxs['text'] = $(this).find('.contentpic').text();
        boxs['color'] = $(this).find('.bg').css('background-color');
        boxs['height'] = $(this).height();
        boxs['width'] = $(this).width();
        boxs['pageX'] = $(this).position().left;
        boxs['pageY'] = $(this).position().top;
        console.dir(boxs);
        data.push(boxs);
    });
});
//创建拖拽方法
function draginit(e) {
    console.log(e)
    e.mousedown(function (e) {
        // var dragroom = $(this);
        // e.preventDefault();
        // var pos = $(this).position();
        // this.posix = { 'x': e.pageX - pos.left, 'y': e.pageY - pos.top };
        // $.extend(document, {
        //     'move': true, 'move_target': this, 'call_down': function (e, posix) {
        //         dragroom.css({
        //             'cursor': 'move',
        //             'top': e.pageY - posix.y,
        //             'left': e.pageX - posix.x
        //         });
        //     }, 'call_up': function () {
        //         dragroom.css('cursor', 'default');
        //     }
        // });
    }).on('mousedown', '.boxs', function (e) {
        if (lock) return;
        var pos = $(this).position();
        this.posix = { 'x': e.pageX - pos.left, 'y': e.pageY - pos.top };
        $.extend(document, { 'move': true, 'move_target': this });
        e.stopPropagation();
    }).on('mousedown', '.boxs .coors', function (e) {
        var $boxs = $(this).parent();
        var posix = {
            'w': $boxs.width(),
            'h': $boxs.height(),
            'x': e.pageX,
            'y': e.pageY
        };
        $.extend(document, {
            'move': true, 'call_down': function (e) {
                $boxs.css({
                    'width': Math.max(30, parseInt((e.pageX - posix.x + posix.w) / 10) * 10),
                    'height': Math.max(40, parseInt((e.pageY - posix.y + posix.h) / 10) * 10)
                });
                var text_w = $($boxs.find('.contentpic')).width()
                var text_h = $($boxs.find('.contentpic')).height()
                text_w = $boxs.find('.contentpic').parent().width() / 2 - text_w / 2
                $boxs.find('.contentpic').css({
                    marginLeft: text_w
                });
            }
        });
        e.stopPropagation();
    });
}
    //测试加载
    // var loadData = [{ id: 1001, text: "C16\n16.5", color: "rgb(255, 0, 0)", height: 70, width: 77, pageX: 627, pageY: 364 },
    // { id: 1002, text: "C17\n16.18", color: "rgb(255, 255, 0)", height: 70, width: 77, pageX: 709, pageY: 364 },
    // { id: 1003, text: "C18\n16.08", color: "rgb(128, 0, 128)", height: 70, width: 77, pageX: 790, pageY: 364 },
    // { id: 1004, text: "C19\n16.08", color: "rgb(0, 128, 0)", height: 70, width: 77, pageX: 870, pageY: 364 },
    // { id: 1005, text: "C20\n16.5", color: "rgb(0, 0, 255)", height: 70, width: 77, pageX: 627, pageY: 439 },
    // { id: 1006, text: "C21\n16.18", color: "rgb(255, 165, 0)", height: 70, width: 77, pageX: 709, pageY: 439 },
    // { id: 1007, text: "C22\n16.08", color: "rgb(255, 165, 0)", height: 70, width: 77, pageX: 870, pageY: 439 },
    // { id: 1008, text: "C23\n16.08", color: "rgb(255, 165, 0)", height: 70, width: 77, pageX: 789, pageY: 439 }];
    // $.each(loadData, function (i, row) {
    //     createBox(row);
    // });
