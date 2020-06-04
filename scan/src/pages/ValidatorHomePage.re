module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let header =
    style([display(`flex), flexDirection(`row), alignItems(`center), height(`px(50))]);

  let validatorsLogo = style([minWidth(`px(50)), marginRight(`px(10))]);
  let highlight = style([margin2(~v=`px(28), ~h=`zero)]);
  let valueContainer = style([display(`flex), justifyContent(`flexStart)]);
  let monikerContainer = style([maxWidth(`px(250))]);

  let emptyContainer =
    style([
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      height(`px(300)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      backgroundColor(white),
      marginBottom(`px(1)),
    ]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.gray7),
    ]);

  let fullWidth =
    style([
      width(`percent(100.0)),
      display(`flex),
      paddingLeft(`px(26)),
      paddingRight(`px(46)),
    ]);

  let sortableTHead = isRight =>
    style([
      display(`flex),
      flexDirection(`row),
      alignItems(`center),
      cursor(`pointer),
      justifyContent(isRight ? `flexEnd : `flexStart),
    ]);

  let sort = style([width(`px(10))]);
  let downIcon = down =>
    style([
      width(`px(8)),
      marginLeft(`pxFloat(1.6)),
      transform(`rotate(`deg(down ? 0. : 180.))),
    ]);
};

module ToggleButton = {
  open Css;

  [@react.component]
  let make = (~isActive, ~setIsActive) => {
    <div className={style([display(`flex), alignItems(`center)])}>
      <div
        onClick={_ => setIsActive(_ => true)}
        className={style([display(`flex), cursor(`pointer)])}>
        <Text value="Active" color=Colors.purple8 />
      </div>
      <HSpacing size=Spacing.sm />
      <div
        className={style([
          display(`flex),
          justifyContent(isActive ? `flexStart : `flexEnd),
          backgroundColor(Colors.gray2),
          borderRadius(`px(15)),
          padding2(~v=`px(2), ~h=`px(3)),
          width(`px(45)),
          cursor(`pointer),
          boxShadow(
            Shadow.box(
              ~inset=true,
              ~x=`zero,
              ~y=`zero,
              ~blur=`px(4),
              isActive ? Colors.purple2 : Colors.gray7,
            ),
          ),
        ])}
        onClick={_ => setIsActive(oldVal => !oldVal)}>
        <img
          src={isActive ? Images.activeValidatorLogo : Images.inactiveValidatorLogo}
          className={style([width(`px(15))])}
        />
      </div>
      <HSpacing size=Spacing.sm />
      <div
        onClick={_ => setIsActive(_ => false)}
        className={style([display(`flex), cursor(`pointer)])}>
        <Text value="Inactive" />
      </div>
    </div>;
  };
};

let renderBody =
    (rank, validatorSub: ApolloHooks.Subscription.variant(ValidatorSub.t), bondedTokenCount) => {
  <TBody
    key={
      switch (validatorSub) {
      | Data({operatorAddress}) => operatorAddress |> Address.toOperatorBech32
      | _ => rank |> string_of_int
      }
    }>
    <div className=Styles.fullWidth>
      <Row>
        <Col size=0.8>
          {switch (validatorSub) {
           | Data(_) =>
             <Text
               value={rank |> string_of_int}
               color=Colors.gray7
               code=true
               weight=Text.Regular
               spacing={Text.Em(0.02)}
               block=true
               size=Text.Md
             />
           | _ => <LoadingCensorBar width=20 height=15 />
           }}
        </Col>
        <Col size=1.9>
          {switch (validatorSub) {
           | Data({operatorAddress, moniker}) =>
             <div className=Styles.monikerContainer>
               <ValidatorMonikerLink validatorAddress=operatorAddress moniker />
             </div>
           | _ => <LoadingCensorBar width=150 height=15 />
           }}
        </Col>
        <Col size=1.4>
          {switch (validatorSub) {
           | Data({tokens, votingPower}) =>
             <div>
               <Text
                 value={tokens |> Coin.getBandAmountFromCoin |> Format.fPretty(~digits=0)}
                 color=Colors.gray7
                 code=true
                 weight=Text.Regular
                 spacing={Text.Em(0.02)}
                 block=true
                 align=Text.Right
                 size=Text.Md
               />
               <VSpacing size=Spacing.sm />
               <Text
                 value={
                   "("
                   ++ (votingPower /. bondedTokenCount *. 100. |> Format.fPercent(~digits=2))
                   ++ ")"
                 }
                 color=Colors.gray6
                 code=true
                 weight=Text.Thin
                 spacing={Text.Em(0.02)}
                 block=true
                 align=Text.Right
                 size=Text.Md
               />
             </div>
           | _ =>
             <>
               <LoadingCensorBar width=100 height=15 isRight=true />
               <VSpacing size=Spacing.sm />
               <LoadingCensorBar width=40 height=15 isRight=true />
             </>
           }}
        </Col>
        <Col size=1.2>
          {switch (validatorSub) {
           | Data({commission}) =>
             <Text
               value={commission->Js.Float.toFixedWithPrecision(~digits=2)}
               color=Colors.gray7
               code=true
               weight=Text.Regular
               spacing={Text.Em(0.02)}
               block=true
               align=Text.Right
               size=Text.Md
             />
           | _ => <LoadingCensorBar width=70 height=15 isRight=true />
           }}
        </Col>
        <Col size=1.1>
          {switch (validatorSub) {
           | Data({uptime}) =>
             <Text
               value={
                 switch (uptime) {
                 | Some(uptime') => uptime'->Js.Float.toFixedWithPrecision(~digits=2)
                 | None => "N/A"
                 }
               }
               color=Colors.gray7
               code=true
               weight=Text.Regular
               spacing={Text.Em(0.02)}
               block=true
               align=Text.Right
               size=Text.Md
             />
           | _ => <LoadingCensorBar width=70 height=15 isRight=true />
           }}
        </Col>
      </Row>
    </div>
  </TBody>;
};

let addUptimeOnValidators =
    (validators: array(ValidatorSub.t), votesBlock: array(ValidatorSub.validator_vote_t)) => {
  validators->Belt.Array.map(validator => {
    let signedBlock =
      votesBlock
      ->Belt.Array.keep(({consensusAddress, voted}) =>
          validator.consensusAddress == consensusAddress && voted == true
        )
      ->Belt.Array.get(0)
      ->Belt.Option.mapWithDefault(0, ({count}) => count)
      |> float_of_int;

    let missedBlock =
      votesBlock
      ->Belt.Array.keep(({consensusAddress, voted}) =>
          validator.consensusAddress == consensusAddress && voted == false
        )
      ->Belt.Array.get(0)
      ->Belt.Option.mapWithDefault(0, ({count}) => count)
      |> float_of_int;

    {
      ...validator,
      uptime:
        signedBlock == 0. && missedBlock == 0.
          ? None : Some(signedBlock /. (signedBlock +. missedBlock) *. 100.),
    };
  });
};

type sort_by_t =
  | NameAsc
  | NameDesc
  | VotingPowerAsc
  | VotingPowerDesc
  | CommissionAsc
  | CommissionDesc
  | UptimeAsc
  | UptimeDesc;

let defaultCompare = (a: ValidatorSub.t, b: ValidatorSub.t) =>
  if (a.tokens != b.tokens) {
    a.tokens.amount < b.tokens.amount ? 1 : (-1);
  } else {
    a.moniker->Js.String.toLocaleLowerCase >= b.moniker->Js.String.toLocaleLowerCase ? 1 : (-1);
  };

let compare = (defaultCompare, condition, execute) =>
  if (condition) {
    execute;
  } else {
    defaultCompare();
  };

let sorting = (validators: array(ValidatorSub.t), sortedBy) => {
  validators
  ->Belt.List.fromArray
  ->Belt.List.sort((a, b) => {
      let curriedCompare = (condition, execute) =>
        compare(() => defaultCompare(a, b), condition, execute);
      switch (sortedBy) {
      | NameAsc =>
        curriedCompare(
          a.moniker != b.moniker,
          a.moniker->Js.String.toLocaleLowerCase > b.moniker->Js.String.toLocaleLowerCase
            ? 1 : (-1),
        )
      | NameDesc =>
        curriedCompare(
          a.moniker != b.moniker,
          a.moniker->Js.String.toLocaleLowerCase < b.moniker->Js.String.toLocaleLowerCase
            ? 1 : (-1),
        )
      | VotingPowerAsc =>
        if (a.tokens != b.tokens) {
          a.tokens.amount > b.tokens.amount ? 1 : (-1);
        } else {
          defaultCompare(a, b);
        }
      | VotingPowerDesc => defaultCompare(a, b)
      | CommissionAsc =>
        if (a.commission != b.commission) {
          a.commission > b.commission ? 1 : (-1);
        } else {
          defaultCompare(a, b);
        }
      | CommissionDesc =>
        if (a.commission != b.commission) {
          a.commission < b.commission ? 1 : (-1);
        } else {
          defaultCompare(a, b);
        }
      | UptimeAsc =>
        a.uptime->Belt.Option.getWithDefault(0.) >= b.uptime->Belt.Option.getWithDefault(0.)
          ? 1 : (-1)
      | UptimeDesc =>
        a.uptime->Belt.Option.getWithDefault(0.) < b.uptime->Belt.Option.getWithDefault(0.)
          ? 1 : (-1)
      };
    })
  ->Belt.List.toArray;
};

module SortableTHead = {
  [@react.component]
  let make = (~title, ~asc, ~desc, ~toggle, ~sortedBy, ~isRight=true) => {
    <div className={Styles.sortableTHead(isRight)} onClick={_ => toggle(asc, desc)}>
      <Text
        block=true
        value=title
        size=Text.Sm
        weight=Text.Semibold
        color=Colors.gray6
        spacing={Text.Em(0.1)}
      />
      <HSpacing size=Spacing.xs />
      {if (sortedBy == asc) {
         <img src=Images.sortDown className={Styles.downIcon(false)} />;
       } else if (sortedBy == desc) {
         <img src=Images.sortDown className={Styles.downIcon(true)} />;
       } else {
         <img src=Images.sort className=Styles.sort />;
       }}
    </div>;
  };
};

module ValidatorList = {
  [@react.component]
  let make = (~allSub) => {
    let (sortedBy, setSortedBy) = React.useState(_ => VotingPowerDesc);

    let toggle = (sortedByAsc, sortedByDesc) =>
      if (sortedBy == sortedByDesc) {
        setSortedBy(_ => sortedByAsc);
      } else {
        setSortedBy(_ => sortedByDesc);
      };

    <>
      <THead>
        <div className=Styles.fullWidth>
          <Row>
            <Col size=0.8 key="RANK">
              <Text
                block=true
                value="RANK"
                size=Text.Sm
                weight=Text.Semibold
                color=Colors.gray6
                spacing={Text.Em(0.1)}
              />
            </Col>
            <Col size=1.9>
              <SortableTHead
                title="VALIDATOR"
                asc=NameAsc
                desc=NameDesc
                toggle
                sortedBy
                isRight=false
              />
            </Col>
            <Col size=1.4>
              <SortableTHead
                title="VOTING POWER (BAND)"
                asc=VotingPowerAsc
                desc=VotingPowerDesc
                toggle
                sortedBy
              />
            </Col>
            <Col size=1.2>
              <SortableTHead
                title="COMMISSION (%)"
                asc=CommissionAsc
                desc=CommissionDesc
                toggle
                sortedBy
              />
            </Col>
            <Col size=1.1>
              <SortableTHead title="UPTIME (%)" asc=UptimeAsc desc=UptimeDesc toggle sortedBy />
            </Col>
          </Row>
        </div>
      </THead>
      {switch (allSub) {
       | ApolloHooks.Subscription.Data((
           (_, _, bondedTokenCount: Coin.t, _, _),
           rawValidators,
           votesBlock,
         )) =>
         let validators = addUptimeOnValidators(rawValidators, votesBlock);
         <>
           {validators->Belt_Array.size > 0
              ? validators
                ->sorting(sortedBy)
                ->Belt_Array.map(e =>
                    renderBody(e.rank, Sub.resolve(e), bondedTokenCount.amount)
                  )
                ->React.array
              : <div className=Styles.emptyContainer>
                  <Text value="No Validators" size=Text.Xxl />
                </div>}
           <VSpacing size=Spacing.lg />
         </>;
       | _ =>
         Belt_Array.make(10, ApolloHooks.Subscription.NoData)
         ->Belt_Array.mapWithIndex((i, noData) => renderBody(i, noData, 1.0))
         ->React.array
       }}
    </>;
  };
};

let getPrevDay = _ => {
  (
    MomentRe.momentNow()
    |> MomentRe.Moment.subtract(~duration=MomentRe.duration(1., `days))
    |> MomentRe.Moment.toUnix
    |> float_of_int
  )
  *. 1000.;
};

let getCurrentDay = _ => {
  (MomentRe.momentNow() |> MomentRe.Moment.toUnix |> float_of_int) *. 1000.;
};

[@react.component]
let make = () => {
  let (prevDayTime, setPrevDayTime) = React.useState(getPrevDay);
  let (currentTime, setCurrentTime) = React.useState(getCurrentDay);

  React.useEffect0(() => {
    let timeOutID =
      Js.Global.setInterval(
        () => {
          setPrevDayTime(getPrevDay);
          setCurrentTime(getCurrentDay);
        },
        60_000,
      );
    Some(() => {Js.Global.clearInterval(timeOutID)});
  });

  let (isActive, setIsActive) = React.useState(_ => true);

  let validatorsSub = ValidatorSub.getList(~isActive, ());
  let validatorsCountSub = ValidatorSub.count();
  let isActiveValidatorCountSub = ValidatorSub.countByActive(isActive);
  let bondedTokenCountSub = ValidatorSub.getTotalBondedAmount();
  let avgBlockTimeSub = BlockSub.getAvgBlockTime(prevDayTime, currentTime);
  let metadataSub = MetadataSub.use();
  let votesBlockSub = ValidatorSub.getListVotesBlock();

  let topPartAllSub =
    Sub.all5(
      validatorsCountSub,
      isActiveValidatorCountSub,
      bondedTokenCountSub,
      avgBlockTimeSub,
      metadataSub,
    );

  let allSub = Sub.all3(topPartAllSub, validatorsSub, votesBlockSub);

  <>
    <Row justify=Row.Between>
      <div className=Styles.header>
        <img src=Images.validators className=Styles.validatorsLogo />
        <Text
          value="ALL VALIDATORS"
          weight=Text.Medium
          size=Text.Md
          nowrap=true
          color=Colors.gray7
          spacing={Text.Em(0.06)}
        />
        {switch (topPartAllSub) {
         | Data((validatorCount, _, _, _, _)) =>
           <>
             <div className=Styles.seperatedLine />
             <Text value={(validatorCount |> string_of_int) ++ " In total"} />
           </>
         | _ => React.null
         }}
      </div>
      <Col>
        {switch (topPartAllSub) {
         | Data(_) => <ToggleButton isActive setIsActive />
         | _ => React.null
         }}
      </Col>
    </Row>
    <div className=Styles.highlight>
      <Row>
        <Col size=0.7>
          {switch (topPartAllSub) {
           | Data((validatorCount, isActiveValidatorCount, _, _, _)) =>
             <InfoHL
               info={InfoHL.Fraction(isActiveValidatorCount, validatorCount, false)}
               header="VALIDATORS"
             />
           | _ =>
             <>
               <LoadingCensorBar width=105 height=15 />
               <VSpacing size=Spacing.sm />
               <LoadingCensorBar width=45 height=15 />
             </>
           }}
        </Col>
        <Col size=1.1>
          {switch (topPartAllSub) {
           | Data((_, _, bondedTokenCount, _, metadata)) =>
             <InfoHL
               info={
                 InfoHL.Fraction(
                   bondedTokenCount |> Coin.getBandAmountFromCoin |> int_of_float,
                   metadata.totalSupply->Coin.getBandAmountFromCoins |> int_of_float,
                   true,
                 )
               }
               header="BONDED TOKENS"
             />
           | _ =>
             <>
               <LoadingCensorBar width=105 height=15 />
               <VSpacing size=Spacing.sm />
               <LoadingCensorBar width=45 height=15 />
             </>
           }}
        </Col>
        <Col size=0.9>
          {switch (topPartAllSub) {
           | Data((_, _, _, _, metadata)) =>
             <InfoHL
               info={InfoHL.FloatWithSuffix(metadata.inflationRate *. 100., "  %", 2)}
               header="INFLATION RATE"
             />
           | _ =>
             <>
               <LoadingCensorBar width=105 height=15 />
               <VSpacing size=Spacing.sm />
               <LoadingCensorBar width=45 height=15 />
             </>
           }}
        </Col>
        <Col size=0.51>
          {switch (topPartAllSub) {
           | Data((_, _, _, avgBlockTime, _)) =>
             <InfoHL
               info={InfoHL.FloatWithSuffix(avgBlockTime, "  secs", 2)}
               header="24 HOUR AVG BLOCK TIME"
             />
           | _ =>
             <>
               <LoadingCensorBar width=105 height=15 />
               <VSpacing size=Spacing.sm />
               <LoadingCensorBar width=45 height=15 />
             </>
           }}
        </Col>
      </Row>
    </div>
    <ValidatorList allSub />
    <VSpacing size=Spacing.lg />
  </>;
};
