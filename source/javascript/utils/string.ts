interface StringConstructor {
    format(msg: string, ...args: any[]): string;
    isNullOrWhiteSpace(s: string): boolean
}

if (!String.format) {
    String.format = function(format: string, ...ar: any[]) {
        var args = Array.prototype.slice.call(arguments, 1);
        return format.replace(/{(\d+)}/g, function(match: any, number: any) {
            return typeof args[number] != 'undefined'
                ? args[number]
                : match
                ;
        });
    };
}

if(!String.isNullOrWhiteSpace) {
    String.isNullOrWhiteSpace = function(str: string){
        return str === null || str.match(/^ *$/) !== null;
    }
}
