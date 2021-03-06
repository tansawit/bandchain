[@react.component]
let make = (~validator, ~setMsgsOpt) => {
  let (srcValidator, setSrcValidator) = React.useState(_ => EnhanceTxInput.empty);
  let (dstValidator, setDstValidator) = React.useState(_ => EnhanceTxInput.empty);
  let (amount, setAmount) = React.useState(_ => EnhanceTxInput.empty);

  React.useEffect3(
    _ => {
      let msgsOpt = {
        let%Opt srcValidatorValue = srcValidator.value;
        let%Opt dstValidatorValue = dstValidator.value;
        let%Opt amountValue = amount.value;
        Some([|
          TxCreator.Redelegate(
            srcValidatorValue,
            dstValidatorValue,
            {amount: amountValue |> Js.Float.toString, denom: "uband"},
          ),
        |]);
      };
      setMsgsOpt(_ => msgsOpt);
      None;
    },
    (srcValidator, dstValidator, amount),
  );

  <>
    <EnhanceTxInput
      width=360
      inputData=srcValidator
      setInputData=setSrcValidator
      parse=Address.fromBech32Opt
      msg="From"
      errMsg="Invalid Address"
      code=true
    />
    <VSpacing size=Spacing.md />
    <EnhanceTxInput
      width=360
      inputData=dstValidator
      setInputData=setDstValidator
      parse=Address.fromBech32Opt
      msg="To"
      errMsg="Invalid Address"
      code=true
    />
    <VSpacing size=Spacing.md />
    <EnhanceTxInput
      width=115
      inputData=amount
      setInputData=setAmount
      parse=Parse.getBandAmount
      msg="Amount (BAND)"
      errMsg="Invalid amount"
      code=true
    />
  </>;
};
