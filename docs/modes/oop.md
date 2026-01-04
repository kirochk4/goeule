# Object Oriented Style

```js
class Archer extends Hero {
    static className = "archer";

    _hp = 80;

    new init(name, attackRange) {
        super.init(name);
        this.attackRange = attackRange;
    }

    sayClass() {
        println(this.className);
    }

    sayHP() {
        println(this._hp);
    }
}

var l = Archer.init("Legolas", 100);

/* == is equal to =========================================================== */

var Archer = Hero {
    className: "archer",
    init() {
        var instance = this {
            ...getPrototypeOf(this).init(name),
            ...{ _hp: 80 },
        };
        instance.attackRange = attackRange;
        return instance;
    },
    sayClass() {
        println(this.className);
    },
    sayHP() {
        println(this._hp);
    },
};

var l = Archer.init("Legolas", 100);
```
