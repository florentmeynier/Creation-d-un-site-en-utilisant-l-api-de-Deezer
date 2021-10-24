open Evaluation
open Term
open Types
open Unification

exception EmptyList
exception NotAList
exception TermCannotBeGenerated
exception GenError

let print_typage_res (t : typage_res) : (string) =
  match t with
  | (lt, st, Fini) -> print_lterm lt ^ " typable avec " ^ print_stype st
  | (lt, st, Echec s) -> s
  | _ -> "Error"

let rec trouve_but (eqs : tequa) : (stype) =
  match eqs with
  | [] -> raise VarNotFound
  | ((VarS v, t)::q) when v = "but" -> t
  | ((t, VarS v)::q) when v = "but" -> t
  | (h::q) -> trouve_but q

let rec typeur (e : env) (lt : lterm) : (typage_res) =
  let eq = gen_equas e lt (VarS "but") in
  match unification eq with
  | (res, Fini) -> (lt, trouve_but res, Fini)
  | (res, Echec s) -> (lt, Nat, Echec((print_lterm lt) ^ " pas typable : " ^ s)) 
  | (res, _) -> (lt, Nat, Echec "Shouldn't happen")
and gen_equas (e : env) (lt : lterm) (st : stype) : (tequa) = 
  match lt with
  | Var v -> [(st, cherche_type v e)]
  | App(Fix, Abs(x, lt2)) -> 
      let nv = fresh_var() in
      let e1 = change_string e x nv in
      let resfix = gen_equas e1 lt2 (VarS nv) in
      (st, (VarS nv))::resfix 
  | App(lt1, lt2) ->
      let nv = fresh_var() in
      gen_equas e lt1 (AppS (VarS nv, st)) @ gen_equas e lt2 (VarS nv) 
  | Abs(v1, lt1) -> 
      let nv1 = fresh_var()
      and nv2 = fresh_var() in
      (st, AppS(VarS nv1, VarS nv2))::(gen_equas ((v1, VarS nv1)::e) lt1 (VarS nv2)) 
  | N n -> [(st, Nat)]
  | Add(lt1, lt2) -> (st, Nat)::(gen_equas e lt1 Nat@gen_equas e lt2 Nat)
  | Sub(lt1, lt2) -> (st, Nat)::(gen_equas e lt1 Nat@gen_equas e lt2 Nat)
  | List l -> 
      (match l with
       | Cons(lt1, q) ->
           let nv = VarS (fresh_var()) in 
           (st, ListS nv)::gen_equas e lt1 nv@gen_equas e (List q) (ListS nv)
       | Nil -> []
      )
  | Head l -> 
    let nv = fresh_var() in
    (st, ForAll(nv, AppS(ListS(VarS nv), VarS nv)))::gen_equas e l (ListS (VarS nv))
  | Tail l ->
    let nv = fresh_var () in
    (st, ForAll(nv, AppS(ListS(VarS nv), VarS nv)))::gen_equas e l (ListS (VarS nv))
  | IfZero(corps, fpos , apos) ->
      let nv = VarS (fresh_var()) in
      let e1 = gen_equas e corps Nat in
      let e2 = gen_equas e fpos nv in
      let e3 = gen_equas e apos nv in
      (st, nv)::e1@e2@e3
  | IfEmpty(corps, fpos, apos) ->
      let nv1 = VarS (fresh_var()) in 
      let nv2 = VarS (fresh_var()) in
      let e1 = gen_equas e corps (ListS nv2) in 
      let e2 = gen_equas e fpos nv1 in
      let e3 = gen_equas e apos nv1 in
      (st, nv1)::e1@e2@e3 
  | Fix -> raise TermCannotBeGenerated
  | Let(v, fpos, apos) -> 
      (match typeur e fpos with
       | (_, res, Fini) -> 
           let nt = generalise e res in
           let ne = modify_env e v nt in
           gen_equas ne apos st
       | _ -> raise GenError
      )
and getStype (lt : lterm) : (stype) =
  match lt with 
  | Var v -> VarS v
  | App(lt1, lt2) -> AppS(getStype lt1, getStype lt2)
  | Abs(v, lt1) -> AppS(VarS v, getStype lt1)
  | N n -> Nat
  | Add(lt1, lt2) -> AppS(Nat, AppS(Nat, Nat))
  | Sub(lt1, lt2) -> AppS(Nat, AppS(Nat, Nat))
  | List l -> 
      (match l with
       | Cons(e, q) -> getStype e
       | Nil -> VarS (fresh_var())
      )
  | _ -> VarS "0"
and change_string (e : env) (str1 : string) (str2 : string) : (env) =
  match e with
  | [] -> []
  | (((s, st)::q)) when s = str1 -> (str1, st)::q
  | ((h::q)) -> h::change_string q str1 str2
and modify_env (e : env) (str : string) (st : stype) : (env) =
  match e with
  | [] -> []
  | ((v, st1)::q) when str = v -> (v, st)::q
  | (h::q) -> h::modify_env q str st
                                
                
let ex_id = (Abs ("x", Var "x"));;
let ex_k = (Abs ("x", Abs ("y", Var "x")));;
let ex_s = (Abs ("x", Abs ("y", Abs ("z", App (App (Var "x", Var "z"), App (Var "y", Var "z"))))));;
let ex_nat1 = (App (Abs ("x", Add(Var "x", N 1)), N 3));;
let ex_omega = (App (Abs ("x", App (Var "x", Var "x")), Abs ("y", App (Var "y", Var "y"))));; 
(*let ex_nat3 = App(ex_nat2, ex_id)*)
let ex_list1 = List(Cons(N 1, Cons(N 2, Cons(N 3, Nil))))
let ex_head1 = Head ex_list1
let ex_izte1 = IfZero(N 0, ex_id, (Abs("y", Var "y")))
let ex_izte2 = IfZero(N 0, ex_id, ex_k)

let t_id = typeur [] ex_id
let t_k = typeur [] ex_k
let t_s = typeur [] ex_s
let t_nat1 = typeur [] ex_nat1
let t_omega = typeur [] ex_omega
    (*let t_nat3 = typeur ex_nat3*)
let t_list1 = typeur [] ex_list1
let t_head1 = typeur [] ex_head1
let t_izte1 = typeur [] ex_izte1
let t_izte2 = typeur [] ex_izte2
;;
print_endline "====================== ex_id";
print_endline (print_typage_res t_id); 
print_endline "====================== ex_k";
print_endline (print_typage_res t_k);;
print_endline "====================== ex_s";
print_endline (print_typage_res t_s);;
print_endline "====================== ex_nat1";
print_endline (print_typage_res t_nat1);;
print_endline "====================== ex_omega";
print_endline (print_typage_res t_omega);; 
print_endline "======================";
(*print_endline (print_typage_res t_nat3);;*)
print_endline "====================== ex_list1";
print_endline (print_tequa (gen_equas [] ex_list1 (VarS "but")));
print_endline (print_typage_res t_list1);;
print_endline "====================== ex_head1";
print_endline (print_tequa (gen_equas [] ex_head1 (VarS "but")));
print_endline (print_typage_res t_head1);;
print_endline "====================== ex_izte1";
print_endline (print_tequa (gen_equas [] ex_izte1 (VarS "but")));
print_endline (print_typage_res t_izte1);;
print_endline "====================== ex_izte2";
print_endline (print_tequa (gen_equas [] ex_izte2 (VarS "but")));
print_endline (print_typage_res t_izte2);;


let list = List(Cons(N 1, Cons(N 2, Cons(N 3, Nil))));;
print_lterm list;;
let eList = gen_equas [] list (VarS "but")
let subList = substitue_partout "T26" (VarS "X") eList
let head = Head list;;
gen_equas [] head (VarS "x");;
print_lterm head
;; 

let corps = N 1
let fpos = N 2
let apos = N 3
let ifzero = IfZero(corps, fpos, apos);;
gen_equas [] ifzero (VarS "x")
;; 
unification_etape (Nat, Nat) [(Nat, Nat)];;
gen_equas [] ex_nat1 (VarS "but");;
unification_etape (VarS "but", AppS (Nat, AppS (Nat, Nat))) [(VarS "but", AppS (Nat, AppS (Nat, Nat)))] 
  
let ex_let = Let("x", N 1, N 2);;
gen_equas [] ex_let (VarS "but");;

ltrcbv (capp (capp (csub (cvar "x") (cvar "y")) (cint 1)) (cint 2));;
ltrcbv (cizte (cint 1) (cvar "x") (cvar "y"));;
