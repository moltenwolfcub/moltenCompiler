$$
\begin{align}
	
	[\textcolor{red}{prog}] &\to [\textcolor{lime}{stmt}]^*
	\\
	[\textcolor{red}{stmt}] &\to \begin{cases}
		\textcolor{cyan}{exit}([\textcolor{lime}{expr}]);\\
		\textcolor{cyan}{var}\space\textcolor{yellow}{name};\\
		\textcolor{yellow}{name}=[\textcolor{lime}{expr}];
	\end{cases}
	\\
	[\textcolor{red}{expr}] &\to \text{intLiteral}

\end{align}
$$


### Tmp:

#### Variable syntax
access var name type;\
private var score int;\
score = 5;

score := 5 (assume private when using :=)
