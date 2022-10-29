# aes

golang aes 没有提供不安全的 EBC 算法
所以如果为了兼容老项目，使用了 EBC 算法，就要用到第三方库。

github.com/andreburgaud/crypt2go

这个库提供了 EBC 算法。
