/**
 * Created by mzj on 2018/7/3.
 *
 * 
 */

/**
 * 当前位置
 * @returns {{x: *, y: *}}
 * @constructor
 */
function getPageScroll() {
    var x, y;
    if (window.pageYOffset) {    // all except IE
        y = window.pageYOffset;
        x = window.pageXOffset;
    } else if (document.documentElement && document.documentElement.scrollTop) {    // IE 6 Strict
        y = document.documentElement.scrollTop;
        x = document.documentElement.scrollLeft;
    } else if (document.body) {    // all other IE
        y = document.body.scrollTop;
        x = document.body.scrollLeft;
    }
    return {x: x, y: y};
}

/**
 * 获取指定元素的 scroll 位置
 * @param id
 * @returns {{x: *, y: *}}
 * @constructor
 */
function getPageScrollByID(id) {
    var curleft = 0, curtop = 0, obj = document.getElementById(id);
    if (obj.offsetParent) {
        curleft = obj.offsetLeft;
        curtop = obj.offsetTop;
        while (obj = obj.offsetParent) {
            curleft += obj.offsetLeft;
            curtop += obj.offsetTop - 60;
        }
    }
    return {x: curleft, y: curtop};
}


/**
 *
 * @param position
 * position: { x: , y:  }
 *
 * @param achorTarget String
 * is changed url anchor after scrolled. This is the target dom ID
 *
 * @param duration
 * 可选, 动画执行时间, 单位是毫秒, 默认值是 100
 *
 * @param interval
 * 可选, 动画间隔, 单位是毫秒, 默认值是 1
 */
window.scrollToSmoothly = function (position) {
    var achorTarget = arguments[1];
    var duration = isNaN(arguments[2]) ? 100 : arguments[1];
    var interval = isNaN(arguments[3]) ? 1 : arguments[2];

    if (isNaN(position.x) || isNaN(position.y)) {
        return ;
    }

    var startPosition = getPageScroll();
    var start_x = startPosition.x,
        start_y = startPosition.y,
        target_x = position.x,
        target_y = position.y,
        variation_x = (target_x - start_x) / duration,
        variation_y = (target_y - start_y) / duration;


    var count_max = duration / interval,
        count = 0;
    var T = setInterval(function() {
        count++;
        if (count < count_max) {
            window.scrollTo(start_x + variation_x * count, start_y + variation_y * count);
        } else {
            if (achorTarget) {
                window.location.href = urlRemoveAnchor(window.location.href) + '#' + achorTarget;
            }
            clearInterval(T);
        }
    }, interval);
};

function urlRemoveAnchor(url) {
    var pos = url.indexOf('#');
    if (pos < 0) pos = url.length;
    return url.slice(0, pos);
}
