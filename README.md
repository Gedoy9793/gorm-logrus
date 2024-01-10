# gorm-logrus

A wrapper between logrus Logger and gorm Logger.

logrus日志到gorm日志的兼容层。


## Example 示例

```go
package main

import (
    "github.com/gedoy9793/gorm-logrus"
    "github.com/sirupsen/logrus"
    "gorm.io/gorm"
	"gorm.io/driver/mysql"
)

func main() {
    logrusLogger := logrus.New()
    
    // you can config this logger
    
    db, _ := gorm.Open(mysql.Open("data.db"), &gorm.Config{
        Logger: gorm_logrus.New(logrusLogger),
    })
}

```