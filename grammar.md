$$
\begin{align}
	
	[\textcolor{red}{prog}] &\to [\textcolor{lime}{stmt}]^*
	\\
	[\textcolor{red}{stmt}] &\to \begin{cases}
		\textcolor{cyan}{exit}([\textcolor{lime}{expr}]);\\
		\textcolor{cyan}{var}\space\textcolor{yellow}{name};\\
		\textcolor{yellow}{name}=[\textcolor{lime}{expr}];\\
		[\textcolor{lime}{scope}]\\
		[\textcolor{lime}{if}]\\
		\textcolor{cyan}{while}([\textcolor{lime}{expr}])[\textcolor{lime}{scope}]\\
		\textcolor{cyan}{break};\\
		\textcolor{cyan}{continue};\\
	\end{cases}
	\\
	[\textcolor{red}{expr}] &\to \begin{cases}
		[\textcolor{lime}{term}]\\
		[\textcolor{lime}{binExpr}]\\
	\end{cases}
	\\
	[\textcolor{red}{binExpr}] &\to \begin{cases}
		[\textcolor{lime}{expr}]*[\textcolor{lime}{expr}] & \textcolor{magenta}{prec=1}\\
		[\textcolor{lime}{expr}]/[\textcolor{lime}{expr}] & \textcolor{magenta}{prec=1}\\
		[\textcolor{lime}{expr}]+[\textcolor{lime}{expr}] & \textcolor{magenta}{prec=0}\\
		[\textcolor{lime}{expr}]-[\textcolor{lime}{expr}] & \textcolor{magenta}{prec=0}\\
	\end{cases}
	\\
	[\textcolor{red}{term}] &\to \begin{cases}
		\text{intLiteral}\\
		\textcolor{yellow}{identifier}\\
		([\textcolor{lime}{expr}])\\
	\end{cases}
	\\
	[\textcolor{red}{scope}] &\to \{[\textcolor{lime}{stmt}]^*\}
	\\
	[\textcolor{red}{if}] &\to \textcolor{cyan}{if}([\textcolor{lime}{expr}])[\textcolor{lime}{scope}]<\textcolor{cyan}{else}\space[\textcolor{lime}{else}]>\\
	[\textcolor{red}{else}] &\to \begin{cases}
		[\textcolor{lime}{if}]\\
		[\textcolor{lime}{scope}]\\
	\end{cases}

\end{align}
$$


### Tmp:

#### Variable syntax
access var name type;\
private var score int;\
score = 5;

score := 5 (assume private when using :=)
