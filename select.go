/**
* @program: esql
*
* @description:
*
* @author: lemo
*
* @create: 2023-08-14 21:41
**/

package esql

import (
	"github.com/xwb1989/sqlparser"
)

func handleSelect(stmt *sqlparser.Select) (dsl string, table string, err error) {

	var query = M{}

	var result = M{
		"query": query,
	}

	var tableName = sqlparser.String(stmt.From)

	// handle where
	var where = stmt.Where
	if where == nil {
		return result.String(), tableName, nil
	}
	handleWhere(query, where.Expr)

	// handle limit
	var limit = stmt.Limit
	if limit != nil {
		result["size"] = Number(limit.Rowcount)
		result["from"] = String(limit.Offset)
	} else {
		result["size"] = 10
		result["from"] = 0
	}

	return result.String(), tableName, nil
}

func handleWhere(result M, expr sqlparser.Expr) {
	switch expr.(type) {
	case *sqlparser.AndExpr:
		handleAnd(result, expr.(*sqlparser.AndExpr))
	case *sqlparser.OrExpr:
		handleOr(result, expr.(*sqlparser.OrExpr))
	case *sqlparser.ParenExpr:
		handleWhere(result, expr.(*sqlparser.ParenExpr).Expr)
	}
}

func handleExpr(result *A, expr sqlparser.Expr, parent sqlparser.Expr) {
	switch expr.(type) {
	case *sqlparser.ComparisonExpr:
		handleComparison(result, expr.(*sqlparser.ComparisonExpr))
	case *sqlparser.IsExpr:
		handleIs(result, expr.(*sqlparser.IsExpr))
	case *sqlparser.RangeCond:
		handleRange(result, expr.(*sqlparser.RangeCond))
	case *sqlparser.AndExpr:
		if _, ok := parent.(*sqlparser.AndExpr); ok {
			handleParentAnd(result, expr.(*sqlparser.AndExpr))
		} else {
			var res = M{}
			*result = append(*result, res)
			handleAnd(res, expr.(*sqlparser.AndExpr))
		}
	case *sqlparser.OrExpr:
		if _, ok := parent.(*sqlparser.OrExpr); ok {
			handleParentOr(result, expr.(*sqlparser.OrExpr))
		} else {
			var res = M{}
			*result = append(*result, res)
			handleOr(res, expr.(*sqlparser.OrExpr))
		}
	case *sqlparser.ParenExpr:
		handleExpr(result, expr.(*sqlparser.ParenExpr).Expr, parent)
	case *sqlparser.FuncExpr:
		// TODO
		panic("not support function")
	default:
		panic("not support " + String(expr))
	}
}

func handleRange(result *A, cond *sqlparser.RangeCond) {

	var field = sqlparser.String(cond.Left)

	var from = FormatSingle(cond.From)

	var to = FormatSingle(cond.To)

	var query = M{
		"range": M{
			field: M{
				"gte": from,
				"lte": to,
			},
		},
	}

	*result = append(*result, query)
}

func handleAnd(result M, expr *sqlparser.AndExpr) {
	var query = &A{}
	result["bool"] = M{"must": query}
	handleExpr(query, expr.Left, expr)
	handleExpr(query, expr.Right, expr)
}

func handleParentAnd(result *A, expr *sqlparser.AndExpr) {
	handleExpr(result, expr.Left, expr)
	handleExpr(result, expr.Right, expr)
}

func handleOr(result M, expr *sqlparser.OrExpr) {
	var query = &A{}
	result["bool"] = M{"should": query}
	handleExpr(query, expr.Left, expr)
	handleExpr(query, expr.Right, expr)
}

func handleParentOr(result *A, expr *sqlparser.OrExpr) {
	handleExpr(result, expr.Left, expr)
	handleExpr(result, expr.Right, expr)
}
