module IsEmpty {
    function isEmpty() {
        var bar:any;
        return function (obj:Object) {
            for (bar in obj) {
                if (obj.hasOwnProperty(bar)) {
                    return false;
                }
            }
            return true;
        };
    }
    angular.module("ModpackHelper").filter('isEmpty', isEmpty);
}
