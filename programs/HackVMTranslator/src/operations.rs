use core::fmt;
use std::str::FromStr;

use enumset::{EnumSetType, EnumSet, enum_set};

/// Operation being performed by a command
#[derive(EnumSetType, Debug)]
pub enum Operation {
    Default,
    /// No operation
    Noop,
    /// Arithmetic operations
    Add,
    Sub,
    Neg,
    Eq,
    Gt,
    Lt,
    And,
    Or,
    Not,
    /// Branching operations
    Label,
    Goto,
    IfGoto,
    /// Memory operations
    Push,
    Pop,
    /// Function operations
    Function,
    Call,
    Return
}

impl FromStr for Operation {
    type Err = ();
    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            ""          => Ok(Operation::Noop),
            "add"       => Ok(Operation::Add),
            "sub"       => Ok(Operation::Sub),
            "neg"       => Ok(Operation::Neg),
            "eq"        => Ok(Operation::Eq),
            "gt"        => Ok(Operation::Gt),
            "lt"        => Ok(Operation::Lt),
            "and"       => Ok(Operation::And),
            "or"        => Ok(Operation::Or),
            "not"       => Ok(Operation::Not),
            "label"     => Ok(Operation::Label),
            "goto"      => Ok(Operation::Goto),
            "if-goto"   => Ok(Operation::IfGoto),
            "push"      => Ok(Operation::Push),
            "pop"       => Ok(Operation::Pop),
            "function"  => Ok(Operation::Function),
            "call"      => Ok(Operation::Call),
            "return"    => Ok(Operation::Return),
            _           => panic!("Read command ({}) is not a supported operation", s)
        }
    }
}

/// Type of operation to be performed
#[derive(EnumSetType, Debug)]
pub enum OperationType {
    Default,
    Arithmetic,
    Branching,
    Memory,
    Function
}

/// The different memory segments supported
#[derive(EnumSetType, Debug)]
pub enum Segment {
    Default,
    None,
    Sp,
    Local,
    Argument,
    This,
    That,
    Constant,
    Static,
    Pointer,
    Temp
}

impl FromStr for Segment {
    type Err = ();
    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            ""          => Ok(Segment::None),
            "sp"        => Ok(Segment::Sp),
            "local"     => Ok(Segment::Local),
            "argument"  => Ok(Segment::Argument),
            "this"      => Ok(Segment::This),
            "that"      => Ok(Segment::That),
            "constant"  => Ok(Segment::Constant),
            "static"    => Ok(Segment::Static),
            "pointer"   => Ok(Segment::Pointer),
            "temp"      => Ok(Segment::Temp),
            _           => panic!("No segment matching {}", s)
        }
    }
}

impl fmt::Display for Segment {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        // Only the below will ever be called in to_string
        match self {
            Segment::None       => write!(f, ""),
            Segment::Sp         => write!(f, "SP"),
            Segment::Local      => write!(f, "LCL"),
            Segment::Argument   => write!(f, "ARG"),
            Segment::This       => write!(f, "THIS"),
            Segment::That       => write!(f, "THAT"),
            _                   => Ok(())
        }
    }
}

/// Set of arithmetic operations
pub const ARITHMETIC_OPERATION: EnumSet<Operation> = enum_set!(
    Operation::Add |
    Operation::Sub |
    Operation::Neg |
    Operation::Eq |
    Operation::Gt |
    Operation::Lt |
    Operation::And |
    Operation::Or |
    Operation::Not
);

/// Set of branching operations
pub const BRANCHING_OPERATION: EnumSet<Operation> = enum_set!(
    Operation::Label |
    Operation::Goto |
    Operation::IfGoto
);

/// Set of memory operations
pub const MEMORY_OPERATION: EnumSet<Operation> = enum_set!(
    Operation::Push |
    Operation::Pop
);

/// Set of function operations
pub const FUNCTION_OPERATION: EnumSet<Operation> = enum_set!(
    Operation::Function |
    Operation::Call |
    Operation::Return
);