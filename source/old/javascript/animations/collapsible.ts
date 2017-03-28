import IAugmentedJQuery = angular.IAugmentedJQuery;
angular.module("ModpackHelper").animation(".collapsible", function () {

    return {
        removeClass: function (elements:IAugmentedJQuery, className: string, done:Function) {
            var tweens: Array<TweenLite> = []
            for(var i = 0; i < elements.length; i++) {
                var element = elements[i];
                var height = element.style.height || 0;
                element.removeAttribute("style");
                var o = TweenLite.from(element, 1, {
                    css: {
                        clearProps: "all",
                        height: height
                    },
                    ease: "Cubic",
                    onComplete: done
                });
                tweens.push(o);
            }
            return function(isCancelled: boolean) {
                if(isCancelled) {
                    tweens.forEach(function(tween) {
                        tween.pause();
                    });
                }
            }
        },
        addClass: function(elements: IAugmentedJQuery, className: string, done:Function) {
            for(var i = 0; i < elements.length; i++) {
                elements[i].removeAttribute("style");
            }
            var o = TweenLite.to(elements, 1, {
                css: {
                    clearProps: "all",
                    height: 0
                },
                ease: "Cubic",
                onComplete: done
            });

            console.log("Leave");

            return function(isCancelled: boolean) {
                if(isCancelled) {
                    o.pause()
                }
            }
        }
    }
});
