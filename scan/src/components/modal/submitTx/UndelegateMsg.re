module Styles = {
  open Css;

  let warning =
    style([
      padding(`px(10)),
      color(Colors.yellow6),
      backgroundColor(Colors.yellow1),
      border(`px(1), `solid, Colors.yellow6),
      borderRadius(`px(4)),
    ]);

  let validator =
    style([display(`flex), flexDirection(`column), alignItems(`flexEnd), width(`px(330))]);

  let info = style([display(`flex), justifyContent(`spaceBetween)]);
};

[@react.component]
let make = (~address, ~validator, ~setMsgsOpt) => {
  let validatorInfoSub = ValidatorSub.get(validator);
  let delegationSub = DelegationSub.getStakeByValiator(address, validator);

  let allSub = Sub.all2(validatorInfoSub, delegationSub);

  let (amount, setAmount) = React.useState(_ => EnhanceTxInput.empty);

  React.useEffect1(
    _ => {
      let msgsOpt = {
        let%Opt amountValue = amount.value;
        Some([|
          TxCreator.Undelegate(
            validator,
            {amount: amountValue |> Js.Float.toString, denom: "uband"},
          ),
        |]);
      };
      setMsgsOpt(_ => msgsOpt);
      None;
    },
    [|amount|],
  );

  <>
    <VSpacing size=Spacing.sm />
    <div className=Styles.warning>
      <Text value="Note: Undelegated balance are locked for 21 days" />
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.info>
      <Text value="Undelegate From" size=Text.Lg spacing={Text.Em(0.03)} nowrap=true block=true />
      {switch (allSub) {
       | Data(({moniker}, _)) =>
         <div className=Styles.validator>
           <Text value=moniker code=true ellipsis=true align=Text.Right />
           <Text
             value={"(" ++ validator->Address.toOperatorBech32 ++ ")"}
             size=Text.Sm
             color=Colors.blueGray5
             code=true
           />
         </div>
       | _ => <LoadingCensorBar width=300 height=26 />
       }}
    </div>
    <VSpacing size=Spacing.lg />
    <VSpacing size=Spacing.md />
    <div className=Styles.info>
      <Text value="Current Stake" size=Text.Lg spacing={Text.Em(0.03)} nowrap=true block=true />
      {switch (allSub) {
       | Data((_, {amount})) =>
         <div>
           <Text
             value={amount |> Coin.getBandAmountFromCoin |> Format.fPretty(~digits=6)}
             code=true
             size=Text.Lg
             weight=Text.Semibold
           />
           <Text value=" BAND" code=true />
         </div>
       | _ => <LoadingCensorBar width=300 height=18 />
       }}
    </div>
    <VSpacing size=Spacing.lg />
    <VSpacing size=Spacing.md />
    <EnhanceTxInput
      width=226
      inputData=amount
      setInputData=setAmount
      parse=Parse.getBandAmount
      msg="Undelegate Amount (BAND)"
      errMsg="Invalid amount"
      placeholder="Insert delegation amount"
      code=true
    />
    <VSpacing size=Spacing.lg />
  </>;
};
