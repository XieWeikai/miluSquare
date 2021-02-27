
$(function () {
    var width = window.innerWidth;
    var height = width * 0.34844437;
    $('.bg-img').css("width",width).css("height",height).css('lineHeight',height+'px');
    //alert(width);
    $('.bg-header').css('width',width).css('bottom',height * 0.05);
    $('body').css('minHeight',window.innerHeight);
    $('footer').css('width',window.innerWidth);
})
window.onresize = function (){
    var width = window.innerWidth;
    var height = width * 0.34834437;
    $('.bg-img').css("width",width).css("height",height).css('lineHeight',height+'px');
    $('.bg-header').css('width',width).css('bottom',height * 0.05);
    $('footer').css('width',window.innerWidth);
}