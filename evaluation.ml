open Term
open Types

type eval_res = lterm * status

let rec ltrcbv_etape (lt : lterm) : (eval_res) =
  match lt with
  | App(fpos, apos) -> 
    (match ltrcbv_etape fpos with
    | (res, Fini) -> (App(res, apos), Fini)
    | _ -> 
      (match ltrcbv_etape apos with
      | (res, Fini) -> (App(fpos, res), Fini)
      | _ ->
        (match (fpos, apos) with
        | (Abs(v, corps), _) -> (instancie corps v apos, Fini)
        | (Head h, List(Cons(e, _))) -> (e, Fini)
        (*| (Tail t, List(Cons(e, q))) -> (List q, Fini)*)
        | (Fix, Abs(v, corps)) -> (instancie (barendregt_lterm corps) v lt, Fini)
        | (App(Add(g, d), N n1), N n2) -> (N (n1 + n2), Fini)
        | (App(Sub(g, d), N n1), N n2) -> (N (n1 - n2), Fini)
        | _ -> (lt, Echec "")
        )
      ) 
    )
  | List l -> ltrcbv_etape_list l
  | Let(v, fpos, apos) -> 
    (match ltrcbv_etape fpos with
    | (res, Fini) -> (Let(v, res, apos), Fini)
    | _ -> (instancie apos v fpos, Fini)
    )
  | IfZero(corps, fpos, apos) ->
    (match ltrcbv_etape corps with
    | (res, Fini) -> (IfZero(res, fpos, apos), Fini)
    | _ -> 
      (match corps with
      | N n when n = 0 -> (fpos, Fini)
      | N _ -> (apos, Fini)
      | _ -> (lt, Echec "")
      )
    )
  | IfEmpty(corps, fpos, apos) ->
    (match ltrcbv_etape corps with
    | (res, Fini) -> (IfEmpty(res, fpos, apos), Fini)
    | _ -> 
      (match corps with
      | List Nil -> (fpos, Fini)
      | List _ -> (apos, Fini)
      | _ -> (lt, Echec "")
      )
    )
  | _ -> (lt, Echec "")
and ltrcbv_etape_list (l : lterm liste) : (eval_res) =
  match l with
  | Nil -> (List Nil, Echec "")
  | Cons(e, q) -> 
    (match ltrcbv_etape e with
    | (res, Fini) -> (List(Cons(res, q)), Fini)
    | _ -> 
      (match ltrcbv_etape_list q with
      | (List res, status) -> (List(Cons(e, res)), status)
      | _ -> (List l, Echec "")
      )
    )

let rec ltrcbv (lt : lterm) : (lterm) =
  let clt = barendregt_lterm lt in
  let res = ltrcbv_rec clt in
  print_endline "------------------------";
  print_endline ("original lterm : " ^ print_lterm lt);
  print_endline ("after barendregt : " ^ print_lterm clt);
  print_endline ("final : " ^ print_lterm res);
  res
and ltrcbv_rec (lt : lterm) : (lterm) =
  match ltrcbv_etape lt with
  | (res, Fini) -> ltrcbv_rec res
  | (res, _) -> res
;;