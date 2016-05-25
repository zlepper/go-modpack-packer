import IAugmentedJQuery = angular.IAugmentedJQuery;
angular.module("ModpackHelper").animation(".collapsible", function () {
    return {
        enter: function (element:IAugmentedJQuery, done:Function) {
            var o = TweenLite.from(element, 1, {
                css: {
                    height: 0
                },
                onComplete: done
            });

            return function(isCancelled: boolean) {
               if(isCancelled) {
                   o.pause();
               }
            }
        }
    }
});
