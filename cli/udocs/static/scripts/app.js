$(document).ready(function() {
    setInitialPopState();
    setSidebarOnPageLoad();
    listenOnSidebarClick();
    listenOnAnchorClick();
    listenOnSearchSubmit();
    listenOnPopstate();
    goToHash();
});

function setInitialPopState() {
    window.popped = (
        'state' in window.history &&
        (typeof window.history.state !== 'undefined') &&
        window.history.state !== null &&
        window.history.state.url != location.href
    )
    window.initialURL = location.href;
}

function setSidebarOnPageLoad() {
    $('#main-sidebar-nav a').each(function() {
        var linkUrl = $(this).attr("href");
        var currentUrl = getLocationPathname(location.pathname) + location.hash;
        if (currentUrl == linkUrl) {
            setSidebar(linkUrl);
        }
    });
}

function listenOnSidebarClick() {
    $(".has-sub-items-content").click(function() {
        $(this).parent(".has-sub-items").toggleClass("is-open");
    });

    $(".sub-items a").click(function(){
        $(".sub-items .is-active").removeClass("is-active");
        $(this).addClass("is-active");
    });
}

function listenOnAnchorClick() {
    $(document).on('click','a', function(event) {
        var url = this.href.toString(),
            path = url.replace(window.location.origin, ""),
            title = '';

        if (isRemoteURL(url.toLowerCase()) || isMediaURL(url.toLowerCase())) {
            return true;
        }

        if (isAnchorTagURL(url.toLowerCase())) {
            title = document.title;
        } else {
            title = $(this).attr('title');
        }

        history.pushState({title: title, path: path}, title, path);
        reportData(title, path);
        return false;
    });
}

function listenOnSearchSubmit() {
    $('#navbar-search').submit(function(event) {
        var path = '/search?q=' + document.getElementById('search-input').value;
        history.pushState({title: 'Search', path: path}, 'Search', path);
        search(path, false);
        return false;
    });
}

function listenOnPopstate() {
    $(window).on('popstate', function (e) {
        var initialPop = !window.popped && location.href == window.initialURL;
        window.popped = true;
        if (!isFirefox && initialPop) {
            history.pushState({title: document.title, path: window.location.pathname}, document.title, window.location.pathname);
            return;
        }
        var state = e.originalEvent.state;
        if (state) {
            if (isSearchURL(state.path)) {
                search(state.path, false);
            } else {
                reportData(state.title, state.path);
            }
            return;
        } else {
            if (isFirefox && isSearchURL(location.pathname)) {
                window.location.reload();
            }
        }
    });
}

function goToHash() {
    if (location.hash) {
        if (navigator.userAgent.indexOf('AppleWebKit') == -1) {
            window.location.hash = location.hash;
        } else {
            window.location.href = location.hash;
            window.location.href = location.hash;
        }
    } else {
        if (!isSearchURL(location.pathname)) {
            var pathname = getLocationPathname(window.location.href);
            if (pathname !== window.location.href) {
                history.pushState({title: 'UDocs', path: pathname}, 'UDocs', pathname);
                window.location.href = pathname;
                window.location.href = pathname;
            }
        }
        $('#main').scrollTop(0);
    }
    expandSidebar();
    if (!isSafari && history.state === null && !isSearchURL(location.pathname)) {
        history.pushState({title: 'UDocs', path: location.pathname}, 'UDocs', location.pathname);
    }
    return false;
};

function isRemoteURL(url) {
    return url.indexOf(window.location.origin.toLowerCase()) == -1;
}

function isAnchorTagURL(url) {
    return url.indexOf('#') >= 0;
}

function isMediaURL(url) {
    var isMedia = false;
    [".png", ".jpeg", ".pptx", ".pdf", ".xml"].forEach(function(suffix){
        if (url.endsWith(suffix)){
            isMedia = true;
        }
    });
    return isMedia;
}

function isSearchURL(url) {
   return url.indexOf('/search') >= 0
}

function redirectSorryPage() {
    if (location.pathname != '/sorry.html') {
        window.setTimeout(function() {
            window.location.href = '/sorry.html'
        }, 300);
    }
    return false;
}

function getLocationPathname(url) {
    if (url.indexOf('.') == -1) {
        if (url.indexOf('/', url.length - 1) !== -1) {
            return location.pathname + 'index.html'
        } else {
            return location.pathname + '/index.html'
        }
    }
    return url;
}

function expandSidebar() {
     var expanded = sessionStorage.getItem('sidebar-expanded');
     var href = getLocationPathname(expanded);
    $('#main-sidebar-nav a[href="' + href + '"]').addClass("is-active").parents('.has-sub-items').addClass("is-open");
}

function setSidebar(href) {
    $('#main-sidebar-nav .is-active').removeClass("is-active");
    $('#main-sidebar-nav a[href="' + href + '"]').addClass("is-active").parents('.has-sub-items').addClass("is-open");
    sessionStorage.setItem('sidebar-expanded', href);
}

function search(url, storeSession) {
    if (storeSession) {
        setSidebar(window.location.pathname);
    }
    window.location = url;
    document.title = 'Search';
}

function reportData(title, path) {
    $.ajax({
        method: "GET",
        url: path.split('#')[0] + '?ajax=true',
        dataType: 'html',
        success: function(res) {
            $('#inner').html(res);
            Prism.highlightAll();
            if (isAnchorTagURL(path)) {
                $('.anchor[name=' + path.split('#')[1] + ']')[0].scrollIntoView();
            } else {
                $('#inner').scrollTop(0);
            }
        }
    });
    document.title = title;
    setSidebar(path.split('#')[0]);
}

var isOpera = (!!window.opr && !!opr.addons) || !!window.opera || navigator.userAgent.indexOf(' OPR/') >= 0;
var isFirefox = typeof InstallTrigger !== 'undefined';
var isSafari = Object.prototype.toString.call(window.HTMLElement).indexOf('Constructor') > 0;
var isIE = /*@cc_on!@*/false || !!document.documentMode;
var isEdge = !isIE && !!window.StyleMedia;
var isChrome = !!window.chrome && !!window.chrome.webstore;
var isBlink = (isChrome || isOpera) && !!window.CSS;
